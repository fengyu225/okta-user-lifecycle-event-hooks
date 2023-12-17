package main

import "fmt"

func GroupLifecycleCreateHandler(event Event) {
	fmt.Println("GroupLifecycleCreate event received.")
	fmt.Println("Group ID: ", event.Target[0].ID)
	fmt.Println("Group Name: ", event.Target[0].DisplayName)
}

func GroupMembershipAddHandler(event Event) {
	fmt.Println("GroupMembershipAdd event received.")
	fmt.Println("Group ID: ", event.Target[0].ID)
	fmt.Println("Group Name: ", event.Target[0].DisplayName)
	fmt.Println("User ID: ", event.Target[1].ID)
	fmt.Println("User Name: ", event.Target[1].DisplayName)
}

func GroupMembershipRemoveHandler(event Event) {
	fmt.Println("GroupMembershipRemove event received.")
	fmt.Println("Group ID: ", event.Target[0].ID)
	fmt.Println("Group Name: ", event.Target[0].DisplayName)
	fmt.Println("User ID: ", event.Target[1].ID)
	fmt.Println("User Name: ", event.Target[1].DisplayName)
}

func GroupProfileUpdateHandler(event Event) {
	fmt.Println("GroupProfileUpdate event received.")
	fmt.Println("Group ID: ", event.Target[0].ID)
	fmt.Println("Group Name: ", event.Target[0].DisplayName)
}

func GroupApplicationAssignmentAddHandler(event Event) {
	fmt.Println("GroupApplicationAssignmentAdd event received.")
	fmt.Println("Group ID: ", event.Target[0].ID)
	fmt.Println("Group Name: ", event.Target[0].DisplayName)
	fmt.Println("Application ID: ", event.Target[1].ID)
	fmt.Println("Application Name: ", event.Target[1].DisplayName)
}

func GroupApplicationAssignmentRemoveHandler(event Event) {
	fmt.Println("GroupApplicationAssignmentRemove event received.")
	fmt.Println("Group ID: ", event.Target[0].ID)
	fmt.Println("Group Name: ", event.Target[0].DisplayName)
	fmt.Println("Application ID: ", event.Target[1].ID)
	fmt.Println("Application Name: ", event.Target[1].DisplayName)
}
