package main

import "fmt"

type autoScaling struct {
	asgName         string
	actionType      string
	dateUTC         string
	capacityMin     string
	capacityDesired string
	capacityMax     string
}

func (p *autoScaling) Register() string {
	cmd := fmt.Sprintf("aws autoscaling put-scheduled-update-group-action --auto-scaling-group-name %s --scheduled-action-name one-time-%s-%s --start-time \"%s\" --min-size %s", p.asgName, p.actionType, p.asgName, p.dateUTC, p.capacityMin)
	if p.capacityMax != "" {
		cmd += fmt.Sprintf(" --max-size %s", p.capacityMax)
	}

	if p.capacityDesired != "" {
		fmt.Printf("aws autoscaling set-desired-capacity --auto-scaling-group-name %s --desired-capacity %s --output json\n", p.asgName, p.capacityDesired)
	}
	cmd += fmt.Sprintln(" --output json")
	return cmd
}

func (p *autoScaling) Deregister() string {
	cmd := fmt.Sprintf("aws autoscaling delete-scheduled-action --auto-scaling-group-name %s --scheduled-action-name one-time-%s-%s", p.asgName, p.actionType, p.asgName)
	cmd += fmt.Sprintln(" --output json")
	return cmd
}
