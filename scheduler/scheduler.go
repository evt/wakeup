package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	scheduler "cloud.google.com/go/scheduler/apiv1"
	schedulerpb "google.golang.org/genproto/googleapis/cloud/scheduler/v1"
)

// Client is a google scheduler API client
type Client struct {
	ctx context.Context
	*scheduler.CloudSchedulerClient
}

// Init creates new scheduler API client instance
func Init(ctx context.Context) (*Client, error) {
	client, err := scheduler.NewCloudSchedulerClient(context.Background())
	if err != nil {
		return nil, err
	}
	return &Client{ctx, client}, nil
}

// CreateJobArgs holds args for CreateJob (see below)
type CreateJobArgs struct {
	RoomNumber        int
	CallTime          string
	CallURL           string
	SchedulerLocation string
	SchedulerTimezone string
}

// CreateJob creates new scheduler job to call provided callURL at provided callTime
// schedulerLocation is a string in format `projects/PROJECT_ID/locations/LOCATION_ID/jobs/JOB_ID`
// * `PROJECT_ID` can contain letters ([A-Za-z]), numbers ([0-9]),
//    hyphens (-), colons (:), or periods (.).
//    For more information, see https://cloud.google.com/resource-manager/docs/creating-managing-projects#identifying_projects
// * `LOCATION_ID` is the canonical ID for the job's location.
//    The list of available locations can be obtained by calling
//    [ListLocations][google.cloud.location.Locations.ListLocations].
//    For more information, see https://cloud.google.com/about/locations/.
// * `JOB_ID` can contain only letters ([A-Za-z]), numbers ([0-9]),
//    hyphens (-), or underscores (_). The maximum length is 500 characters.
func (cl *Client) CreateJob(args CreateJobArgs) error {
	if args.CallTime == "" {
		return fmt.Errorf("[Room %d] no call time provided", args.RoomNumber)
	}
	parts := strings.Split(args.CallTime, ":")
	if len(parts) != 2 {
		return fmt.Errorf("[Room %d] wake up time (%s) must be in the following format: hh:mm", args.RoomNumber, args.CallTime)
	}
	// Parse wake up time hour and min
	callHour, err := strconv.Atoi(parts[0])
	if err != nil {
		return errors.Wrapf(err, "[Room %d] wake up hour (%s) in wake up time (%s) is not a number", args.RoomNumber, parts[0], args.CallTime)
	}
	callMin, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.Wrapf(err, "[Room %d] wake up min (%s) in wake up time (%s) is not a number", args.RoomNumber, parts[1], args.CallTime)
	}
	if args.CallURL == "" {
		return fmt.Errorf("[Room %d] No call URL provided", args.RoomNumber)
	}
	// Prepare scheduler job name
	jobID := args.SchedulerLocation + "/jobs/call_at_" + args.CallTime
	// Make sure scheduler job doesn't exist for provided wake up time
	existingJob, err := cl.GetJob(jobID)
	if err != nil {
		return err
	}
	// Job already exists - do nothing
	if existingJob != nil {
		return nil
	}
	// Prepare schedule to call once a day at provided time
	schedule := fmt.Sprintf("%d %d * * *", callMin, callHour)
	req := &schedulerpb.CreateJobRequest{
		Parent: args.SchedulerLocation,
		Job: &schedulerpb.Job{
			Name: removeColon(jobID),
			Target: &schedulerpb.Job_HttpTarget{
				HttpTarget: &schedulerpb.HttpTarget{
					Uri:        args.CallURL,
					HttpMethod: schedulerpb.HttpMethod_GET,
				},
			},
			TimeZone: args.SchedulerTimezone,
			Schedule: schedule,
		},
	}
	resp, err := cl.CloudSchedulerClient.CreateJob(cl.ctx, req)
	if err != nil {
		return errors.Wrap(err, "CreateJob->cl.CloudSchedulerClient.CreateJob")
	}
	// Make sure job has been created with correct state
	if resp.State != schedulerpb.Job_ENABLED {
		return errors.Wrap(fmt.Errorf("Found incorrect job state %d in response. Must be %d", resp.State, schedulerpb.Job_ENABLED), "CreateJob->cl.CloudSchedulerClient.CreateJob")
	}

	return nil
}

// GetJob returns scheduler job by ID from google cloud scheduler
func (cl *Client) GetJob(jobID string) (*schedulerpb.Job, error) {
	if jobID == "" {
		return nil, errors.New("No job ID provided")
	}
	req := &schedulerpb.GetJobRequest{
		Name: removeColon(jobID),
	}
	resp, err := cl.CloudSchedulerClient.GetJob(cl.ctx, req)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), "or the resource may not exist") {
			return nil, nil
		}
		if strings.Contains(fmt.Sprintf("%s", err), "Job not found") {
			return nil, nil
		}
		return nil, errors.Wrap(err, "GetJob->cl.CloudSchedulerClient.GetJob")
	}
	return resp, nil
}

// removeColon replaces ":" with "_" as it's not allowed in job ID
func removeColon(s string) string {
	return strings.Replace(s, ":", "_", -1)
}
