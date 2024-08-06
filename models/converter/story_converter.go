package converter

import "gopro/models"

func MapStoryToCreateStoryResponse(story *models.Story) *models.CreateStoryResponse {
	return &models.CreateStoryResponse{
		ID:        story.ID,
		Status:    story.Status,
		CreatedAt: story.CreatedAt,
		UpdatedAt: story.UpdatedAt,
	}
}

func MapStoryToUpdateStoryResponse(story *models.Story) *models.UpdateStoryResponse {
	return &models.UpdateStoryResponse{
		ID:        story.ID,
		Status:    story.Status,
		CreatedAt: story.CreatedAt,
		UpdatedAt: story.UpdatedAt,
	}
}

func MapStoryToGetStoriesResponse(story *models.Story) *models.GetStoriesResponse {
	return &models.GetStoriesResponse{
		ID:        story.ID,
		Status:    story.Status,
		UserID:    story.UserID,
		CreatedAt: story.CreatedAt,
		UpdatedAt: story.UpdatedAt,
	}
}

func MapStoriesToGetStoriesResponse(stories *[]models.Story) []*models.GetStoriesResponse {
	var storiesResponse []*models.GetStoriesResponse
	for _, story := range *stories {
		storyResponse := MapStoryToGetStoriesResponse(&story)
		storiesResponse = append(storiesResponse, storyResponse)
	}
	return storiesResponse
}
