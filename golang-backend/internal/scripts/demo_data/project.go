package demodatascripts

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
	"github.com/aportela/doneo/internal/services/projectservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/gofrs/uuid"
)

func getRandomProjectSummary() string {
	actions := []string{
		"Develop a solution",
		"Implement a system",
		"Conduct a market study",
		"Create a marketing strategy",
		"Build a new product",
		"Design a user interface",
		"Optimize existing workflows",
		"Enhance customer service",
	}

	results := []string{
		"to increase efficiency",
		"with the goal of improving customer experience",
		"to optimize internal processes",
		"to expand market presence",
		"to enhance user engagement",
		"to boost sales conversion rates",
		"to improve operational performance",
		"to maximize ROI",
	}

	goals := []string{
		"with a focus on sustainability",
		"to achieve long-term growth",
		"to reduce operational costs",
		"to improve brand awareness",
		"to support business scalability",
		"to foster innovation",
		"to enhance employee productivity",
		"to improve time-to-market",
	}

	action := actions[rand.Intn(len(actions))]
	result := results[rand.Intn(len(results))]
	goal := goals[rand.Intn(len(goals))]
	return fmt.Sprintf("%s %s %s.", action, result, goal)
}

func getRandomProjectDescription() string {
	approaches := []string{
		"Adopt a user-centric approach",
		"Implement agile methodologies",
		"Utilize machine learning to enhance features",
		"Focus on automation and process improvement",
		"Leverage cloud technologies for scalability",
		"Create an intuitive and responsive interface",
		"Optimize existing systems for greater efficiency",
		"Use data-driven insights to make decisions",
	}

	challenges := []string{
		"Handling large volumes of data",
		"Meeting tight deadlines with limited resources",
		"Overcoming integration issues with legacy systems",
		"Addressing scalability challenges",
		"Ensuring security and privacy compliance",
		"Managing stakeholder expectations",
		"Reducing system downtime during migration",
		"Ensuring smooth cross-functional collaboration",
	}

	solutions := []string{
		"By implementing robust data processing pipelines",
		"Using cutting-edge technologies like AI and cloud services",
		"Through a phased rollout and comprehensive testing",
		"By improving collaboration tools and communication channels",
		"By introducing automated workflows and continuous integration",
		"By providing comprehensive training and user support",
		"By leveraging analytics for real-time decision-making",
		"By prioritizing security throughout the project lifecycle",
	}

	approach := approaches[rand.Intn(len(approaches))]
	challenge := challenges[rand.Intn(len(challenges))]
	solution := solutions[rand.Intn(len(solutions))]

	return fmt.Sprintf("Approach: %s.\nChallenge: %s.\nSolution: %s.", approach, challenge, solution)
}

func getRandomProjectKey() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength := 6
	result := make([]byte, keyLength)
	for i := 0; i < keyLength; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func getRandomProject(workspaceId string, userIds []string, projectTypeIds []string) domain.Project {
	projectID := func() string { u, _ := uuid.NewV7(); return u.String() }()
	projectDescription := getRandomProjectDescription()
	startOffset := rand.Int63n(48)
	finishOffset := rand.Int63n(96)
	dueOffset := rand.Int63n(144)
	ctime := utils.GetRandomMSTimestamp(time.Now().AddDate(-5, 0, 0), time.Now())
	stime := ctime + startOffset*int64(time.Hour/time.Millisecond)
	ftime := stime + finishOffset*int64(time.Hour/time.Millisecond)
	dtime := ftime + dueOffset*int64(time.Hour/time.Millisecond)
	rand.Shuffle(len(userIds), func(i, j int) {
		userIds[i], userIds[j] = userIds[j], userIds[i]
	})
	/*
		numParticipants := rand.Intn(3) + 1
		var participants []domain.UserBase
		for i := 0; i < numParticipants; i++ {
			participants = append(participants, domain.UserBase{ID: userIds[i]})
		}
	*/
	return domain.Project{
		ID:          projectID,
		WorkspaceId: workspaceId,
		Key:         getRandomProjectKey(),
		Summary:     getRandomProjectSummary(),
		Description: &projectDescription,
		CreatedBy:   domain.UserBase{ID: userIds[rand.Intn(len(userIds))]},
		CreatedAt:   ctime,
		UpdatedAt:   &ftime,
		StartedAt:   &stime,
		DueAt:       &dtime,
		Type:        domain.ProjectType{ID: projectTypeIds[rand.Intn(len(projectTypeIds))]},
	}
}

func createProjects(database database.Database, workspaceId string, projectTypeIds []string, userIds []string, count int) []string {
	var newProjectIds []string
	projectRepository := projectrepository.NewProjectRepository(database)

	projectService := projectservice.NewProjectService(projectRepository)
	for i := 1; i <= count; i++ {
		newProject := getRandomProject(workspaceId, userIds, projectTypeIds)
		err := projectService.AddProject(context.Background(), newProject)
		if err != nil {
			fmt.Printf("Error creating project %s\n", err.Error())
		}
		newProjectIds = append(newProjectIds, newProject.ID)
	}
	return newProjectIds
}
