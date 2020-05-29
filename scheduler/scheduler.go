package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	scheduler "cloud.google.com/go/scheduler/apiv1"
	"github.com/davecgh/go-spew/spew"
	schedulerpb "google.golang.org/genproto/googleapis/cloud/scheduler/v1"
)

func CreateJob(ctx context.Context, wakeUpTime, callURL, schedulerLocation string) error {
	if wakeUpTime == "" {
		return errors.New("No wake up time provided")
	}
	if len(wakeUpTime) != 5 {
		return fmt.Errorf("Wake up time (%s) must be in the following format: xx:yy", wakeUpTime)
	}
	parts := strings.Split(wakeUpTime, ":")
	if len(parts) != 2 {
		return fmt.Errorf("Wake up time (%s) must be in the following format: xx:yy", wakeUpTime)
	}
	// Parse wake up time hour and min
	wakeUpHour, err := strconv.Atoi(parts[0])
	if err != nil {
		return errors.Wrapf(err, "wake up hour (%s) in wake up time (%s) is not a number", parts[0], wakeUpTime)
	}
	wakeUpMin, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.Wrapf(err, "wake up min (%s) in wake up time (%s) is not a number", parts[1], wakeUpTime)
	}
	if callURL == "" {
		return errors.New("No call URL provided")
	}
	// Create cloud scheduler job
	schedulerClient, err := scheduler.NewCloudSchedulerClient(ctx)
	if err != nil {
		return err
	}
	// Prepare schedule to call once a day at provided time
	schedule := fmt.Sprintf("%d %d * * *", wakeUpMin, wakeUpHour)
	req := &schedulerpb.CreateJobRequest{
		Parent: schedulerLocation,
		Job: &schedulerpb.Job{
			Target: &schedulerpb.Job_HttpTarget{
				HttpTarget: &schedulerpb.HttpTarget{
					Uri:        callURL,
					HttpMethod: schedulerpb.HttpMethod_GET,
				},
			},
			Schedule: schedule,
		},
	}
	resp, err := schedulerClient.CreateJob(ctx, req)
	if err != nil {
		return err
	}
	spew.Dump(resp)
	return nil
}
