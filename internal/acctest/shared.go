package acctest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

var (
	TestTeamName    = "tf-team-shared"
	TestTeam        apiclient.Team
	TestProjectName = "tf-project-shared"
	TestProject     apiclient.Project
)

func SetupShared(ctx context.Context) {
	must.Do(setupSharedTeam(ctx))
	must.Do(setupSharedProject(ctx))
}

func setupSharedTeam(ctx context.Context) error {
	readHttpResp, err := SharedApiClient.GetOrganizationTeamWithResponse(ctx, TestOrganization, TestTeamName)
	if err != nil {
		return err
	} else if readHttpResp.StatusCode() == http.StatusOK && readHttpResp.JSON200 != nil {
		TestTeam = *readHttpResp.JSON200
		return nil
	}

	createHttpResp, err := SharedApiClient.CreateOrganizationTeamWithResponse(ctx, TestOrganization, apiclient.CreateOrganizationTeamJSONRequestBody{
		Name: TestTeamName,
	})
	if err != nil {
		return err
	} else if createHttpResp.StatusCode() == http.StatusCreated && createHttpResp.JSON201 != nil {
		TestTeam = *createHttpResp.JSON201
		return nil
	} else {
		return fmt.Errorf("failed to create shared team: status code=%d, body=%s", createHttpResp.StatusCode(), createHttpResp.Body)
	}
}

func setupSharedProject(ctx context.Context) error {
	readHttpResp, err := SharedApiClient.GetOrganizationProjectWithResponse(ctx, TestOrganization, TestProjectName)
	if err != nil {
		return err
	} else if readHttpResp.StatusCode() == http.StatusOK && readHttpResp.JSON200 != nil {
		TestProject = *readHttpResp.JSON200
		return nil
	}

	createHttpResp, err := SharedApiClient.CreateOrganizationTeamProjectWithResponse(ctx, TestOrganization, TestTeam.Id, apiclient.CreateOrganizationTeamProjectJSONRequestBody{
		Name:     TestProjectName,
		Platform: new("go"),
	})
	if err != nil {
		return err
	} else if createHttpResp.StatusCode() == http.StatusCreated && createHttpResp.JSON201 != nil {
		TestProject = *createHttpResp.JSON201
		return nil
	} else {
		return fmt.Errorf("failed to create shared project: status code=%d, body=%s", createHttpResp.StatusCode(), createHttpResp.Body)
	}
}

func TeardownShared(ctx context.Context) {
	must.Do(teardownSharedTeam(ctx))
	must.Do(teardownSharedProject(ctx))
}

func teardownSharedTeam(ctx context.Context) error {
	httpResp, err := SharedApiClient.DeleteOrganizationTeamWithResponse(ctx, TestOrganization, TestTeam.Id)
	if err != nil {
		return err
	} else if httpResp.StatusCode() != http.StatusNoContent && httpResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to delete shared team: status code=%d, body=%s", httpResp.StatusCode(), httpResp.Body)
	}
	return nil
}

func teardownSharedProject(ctx context.Context) error {
	httpResp, err := SharedApiClient.DeleteOrganizationProjectWithResponse(ctx, TestOrganization, TestProject.Slug)
	if err != nil {
		return err
	} else if httpResp.StatusCode() != http.StatusNoContent && httpResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to delete shared project: status code=%d, body=%s", httpResp.StatusCode(), httpResp.Body)
	}
	return nil
}
