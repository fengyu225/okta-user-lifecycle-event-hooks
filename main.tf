# ================================================
# VPC
# ================================================
resource "aws_vpc" "main" {
  cidr_block           = "192.168.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "Okta-Event-Hooks-Test-VPC"
  }
}

resource "aws_subnet" "public_subnet_us_east_1a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "192.168.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "us-east-1a"
}

resource "aws_subnet" "public_subnet_us_east_1b" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "192.168.2.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "us-east-1b"
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }
}

resource "aws_route_table_association" "public_rt_assoc_a" {
  subnet_id      = aws_subnet.public_subnet_us_east_1a.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_route_table_association" "public_rt_assoc_b" {
  subnet_id      = aws_subnet.public_subnet_us_east_1b.id
  route_table_id = aws_route_table.public_rt.id
}

# ================================================
# Security Group
# ================================================
resource "aws_security_group" "okta_event_hooks_alb" {
  name        = "okta-event-hooks-alb"
  description = "Security group for Okta Event Hooks ALB"
  vpc_id      = aws_vpc.main.id
}

resource "aws_security_group_rule" "okta_event_hooks_alb_egress_all" {
  type        = "egress"
  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.okta_event_hooks_alb.id
}

resource "aws_security_group_rule" "okta_event_hooks_ingress_all" {
  type        = "ingress"
  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.okta_event_hooks_alb.id
}

#resource "aws_security_group_rule" "okta_event_hooks_ingress_okta" {
#  for_each = {for idx, chunk in local.ip_chunks : idx => chunk}
#
#  type              = "ingress"
#  from_port         = 80
#  to_port           = 80
#  protocol          = "tcp"
#  cidr_blocks       = each.value
#  security_group_id = aws_security_group.okta_event_hooks_alb.id
#}

resource "aws_security_group" "okta_event_hooks_ec2" {
  name        = "okta-event-hooks-ec2"
  description = "Security group for Okta Event Hooks EC2"
  vpc_id      = aws_vpc.main.id
}

resource "aws_security_group_rule" "okta_event_hooks_ec2_ingress_ssh" {
  type      = "ingress"
  from_port = 22
  to_port   = 22
  protocol  = "tcp"

  security_group_id = aws_security_group.okta_event_hooks_ec2.id
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "okta_event_hooks_ec2_ingress_alb" {
  type      = "ingress"
  from_port = 8080
  to_port   = 8080
  protocol  = "tcp"

  source_security_group_id = aws_security_group.okta_event_hooks_alb.id
  security_group_id        = aws_security_group.okta_event_hooks_ec2.id
}

resource "aws_security_group_rule" "okta_event_hooks_ec2_egress_all" {
  type        = "egress"
  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.okta_event_hooks_ec2.id
}

# ================================================
# ALB
# ================================================
resource "aws_alb" "okta_event_hooks" {
  name               = "okta-event-hooks-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.okta_event_hooks_alb.id]
  subnets            = [aws_subnet.public_subnet_us_east_1a.id, aws_subnet.public_subnet_us_east_1b.id]
  tags               = {
    Name = "okta-event-hooks-alb"
  }
}

resource "aws_alb_target_group" "okta_event_hooks" {
  name     = "okta-event-hooks-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id
  tags     = {
    Name = "okta-event-hooks-tg"
  }
  health_check {
    path                = "/health"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    matcher             = "200"
  }
}

resource "aws_alb_listener" "okta_event_hooks" {
  load_balancer_arn = aws_alb.okta_event_hooks.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "arn:aws:acm:us-east-1:072422391281:certificate/cf91c5af-0ade-401b-acbe-0e7e330981d9"

  default_action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.okta_event_hooks.arn
  }
}

# ================================================
# Okta Event Hooks ASG
# ================================================
resource "aws_iam_role" "okta_event_hooks" {
  name               = "okta-event-hooks"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ecr_policy" {
  name        = "okta-event-hooks-ecr-policy"
  description = "Policy for allowing access to ECR"

  policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [
      {
        Sid      = "ecrAuthZ"
        Effect   = "Allow"
        Resource = "*"
        Action   = [
          "ecr:GetAuthorizationToken",
        ]
      },
      {
        Sid      = "ecr"
        Effect   = "Allow"
        Resource = "arn:aws:ecr:us-east-1:072422391281:repository/okta-event-hooks"
        Action   = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:GetRepositoryPolicy",
          "ecr:DescribeRepositories",
          "ecr:ListImages",
          "ecr:DescribeImages",
          "ecr:BatchGetImage",
          "ecr:GetLifecyclePolicy",
          "ecr:GetLifecyclePolicyPreview",
          "ecr:ListTagsForResource",
          "ecr:DescribeImageScanFindings"
        ]
      }
    ],
  })
}

resource "aws_iam_role_policy_attachment" "ecr_policy_attachment" {
  role       = aws_iam_role.okta_event_hooks.name
  policy_arn = aws_iam_policy.ecr_policy.arn
}

resource "aws_iam_instance_profile" "okta_event_hooks" {
  name = "okta-event-hooks-instance-profile"
  role = aws_iam_role.okta_event_hooks.name
}

resource "aws_launch_configuration" "okta_event_hooks" {
  name_prefix          = "okta-event-hooks-"
  image_id             = "ami-0fc5d935ebf8bc3bc"
  instance_type        = "t2.medium"
  iam_instance_profile = aws_iam_instance_profile.okta_event_hooks.name
  security_groups      = [aws_security_group.okta_event_hooks_ec2.id]
  key_name             = "yu-feng-uf-1"
  user_data            = file("user-data.sh")
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "okta_event_hooks" {
  name                      = "okta-event-hooks-asg"
  launch_configuration      = aws_launch_configuration.okta_event_hooks.id
  min_size                  = 1
  max_size                  = 1
  desired_capacity          = 1
  vpc_zone_identifier       = [aws_subnet.public_subnet_us_east_1a.id, aws_subnet.public_subnet_us_east_1b.id]
  target_group_arns         = [aws_alb_target_group.okta_event_hooks.arn]
  health_check_type         = "EC2"
  health_check_grace_period = 300
  tag {
    key                 = "Name"
    value               = "okta-event-hooks"
    propagate_at_launch = true
  }
}