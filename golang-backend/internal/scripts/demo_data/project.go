package demodatascripts

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/noterepository"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
	"github.com/aportela/doneo/internal/services/noteservice"
	"github.com/aportela/doneo/internal/services/projectpermissionservice"
	"github.com/aportela/doneo/internal/services/projectservice"
	"github.com/aportela/doneo/internal/utils"
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

func getRandomProject(userIds []string, projectTypeIds []string, projectPriorityIds []string, projectStatusIds []string) domain.Project {
	projectID := utils.UUID()
	projectDescription := getRandomProjectDescription()
	startOffset := rand.Int63n(48)
	finishOffset := rand.Int63n(96)
	dueOffset := rand.Int63n(144)
	ctime := utils.GetRandomMSTimestamp(time.Now().AddDate(-5, 0, 0), time.Now())
	utime := ctime + startOffset*int64(time.Hour/time.Millisecond)
	stime := ctime + startOffset*int64(time.Hour/time.Millisecond)
	ftime := stime + finishOffset*int64(time.Hour/time.Millisecond)
	dtime := ftime + dueOffset*int64(time.Hour/time.Millisecond)
	rand.Shuffle(len(userIds), func(i, j int) {
		userIds[i], userIds[j] = userIds[j], userIds[i]
	})
	rand.Shuffle(len(projectTypeIds), func(i, j int) {
		projectTypeIds[i], projectTypeIds[j] = projectTypeIds[j], projectTypeIds[i]
	})
	rand.Shuffle(len(projectStatusIds), func(i, j int) {
		projectStatusIds[i], projectStatusIds[j] = projectStatusIds[j], projectStatusIds[i]
	})
	return domain.Project{
		ID:          projectID,
		Key:         getRandomProjectKey(),
		Summary:     getRandomProjectSummary(),
		Description: &projectDescription,
		CreatedBy:   domain.UserBase{ID: userIds[rand.Intn(len(userIds))]},
		CreatedAt:   time.UnixMilli(ctime),
		UpdatedAt:   utils.Int64PtrToTimePtr(&utime),
		StartedAt:   utils.Int64PtrToTimePtr(&stime),
		FinishedAt:  utils.Int64PtrToTimePtr(&ftime),
		DueAt:       utils.Int64PtrToTimePtr(&dtime),
		Type:        domain.ProjectType{ID: projectTypeIds[rand.Intn(len(projectTypeIds))]},
		Priority:    domain.ProjectPriority{ID: projectPriorityIds[rand.Intn(len(projectPriorityIds))]},
		Status:      domain.ProjectStatus{ID: projectStatusIds[rand.Intn(len(projectStatusIds))]},
	}
}

func randomText(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz "

	b := make([]byte, n)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
func createProjects(database database.Database, projectTypeIds []string, projectPriorityIds []string, projectStatusIds []string, userIds []string, roleIds []string, count int) []string {
	var newProjectIds []string
	projectService := projectservice.NewService(database, projectrepository.NewRepository(database))
	noteRepository := noterepository.NewRepository(database)
	noteService := noteservice.NewService(noteRepository)
	projectPermissionRepository := projectpermissionrepository.NewRepository(database)
	projectPermissionService := projectpermissionservice.NewService(projectPermissionRepository)
	for i := 1; i <= count; i++ {
		newProject := getRandomProject(userIds, projectTypeIds, projectPriorityIds, projectStatusIds)
		err := projectService.Add(context.Background(), newProject)
		if err != nil {
			fmt.Printf("Error creating project %s\n", err.Error())
		}
		rand.Shuffle(len(userIds), func(i, j int) {
			userIds[i], userIds[j] = userIds[j], userIds[i]
		})
		projectPermissionService.Add(context.Background(), utils.UUID(), newProject.ID, userIds[0], roleIds[0])
		projectPermissionService.Add(context.Background(), utils.UUID(), newProject.ID, userIds[1], roleIds[1])
		projectPermissionService.Add(context.Background(), utils.UUID(), newProject.ID, userIds[2], roleIds[1])
		projectPermissionService.Add(context.Background(), utils.UUID(), newProject.ID, userIds[3], roleIds[1])
		newProjectIds = append(newProjectIds, newProject.ID)
		for j := 0; j < 5; j++ {
			note := domain.Note{
				ID: utils.UUID(),
				User: domain.UserBase{
					ID: userIds[j],
				},
				Body:      "Note index " + strconv.Itoa(j) + ": " + randomText(rand.Intn(384)+128),
				CreatedAt: time.Now().Add(time.Duration(j*5) * time.Minute),
			}
			noteService.AddProjectNote(context.Background(), newProject.ID, note)
		}
	}
	return newProjectIds
}
