package save

import (
	"fmt"
	"log"
	"sort"
)

// UpdateHighScore updates the high score for a specific level
// Returns true if a new high score was set
func (sm *SaveManager) UpdateHighScore(levelID string, score int) bool {
	data := sm.saves[sm.currentSlot]
	if data == nil {
		log.Printf("Warning: no save data in current slot %d", sm.currentSlot)
		return false
	}

	if data.LevelScores == nil {
		data.LevelScores = make(map[string]int)
	}

	currentHigh := data.LevelScores[levelID]
	if score > currentHigh {
		data.LevelScores[levelID] = score
		if err := sm.SaveGame(sm.currentSlot, data); err != nil {
			log.Printf("Error saving high score: %v", err)
			return false
		}
		log.Printf("New high score for %s: %d (previous: %d)", levelID, score, currentHigh)
		return true // New high score
	}

	return false
}

// GetHighScore returns the high score for a specific level
func (sm *SaveManager) GetHighScore(levelID string) int {
	data := sm.saves[sm.currentSlot]
	if data == nil || data.LevelScores == nil {
		return 0
	}
	return data.LevelScores[levelID]
}

// UpdateWorldScore updates the total score for a world
// Returns true if a new high score was set
func (sm *SaveManager) UpdateWorldScore(worldID int, score int) bool {
	data := sm.saves[sm.currentSlot]
	if data == nil {
		log.Printf("Warning: no save data in current slot %d", sm.currentSlot)
		return false
	}

	if data.WorldScores == nil {
		data.WorldScores = make(map[int]int)
	}

	currentHigh := data.WorldScores[worldID]
	if score > currentHigh {
		data.WorldScores[worldID] = score
		if err := sm.SaveGame(sm.currentSlot, data); err != nil {
			log.Printf("Error saving world score: %v", err)
			return false
		}
		log.Printf("New high score for world %d: %d (previous: %d)", worldID, score, currentHigh)
		return true
	}

	return false
}

// GetWorldScore returns the high score for a specific world
func (sm *SaveManager) GetWorldScore(worldID int) int {
	data := sm.saves[sm.currentSlot]
	if data == nil || data.WorldScores == nil {
		return 0
	}
	return data.WorldScores[worldID]
}

// GetTotalScore returns the total score across all levels
func (sm *SaveManager) GetTotalScore() int {
	data := sm.saves[sm.currentSlot]
	if data == nil {
		return 0
	}
	return data.TotalScore
}

// UpdateTotalScore updates the total score
func (sm *SaveManager) UpdateTotalScore(score int) error {
	data := sm.saves[sm.currentSlot]
	if data == nil {
		return fmt.Errorf("no save data in current slot %d", sm.currentSlot)
	}

	data.TotalScore = score
	return sm.SaveGame(sm.currentSlot, data)
}

// AddToTotalScore adds points to the total score
func (sm *SaveManager) AddToTotalScore(points int) error {
	data := sm.saves[sm.currentSlot]
	if data == nil {
		return fmt.Errorf("no save data in current slot %d", sm.currentSlot)
	}

	data.TotalScore += points
	return sm.SaveGame(sm.currentSlot, data)
}

// LevelScoreEntry represents a level score entry for leaderboards
type LevelScoreEntry struct {
	LevelID string
	Score   int
}

// GetTopLevelScores returns the top N level scores across all levels
func (sm *SaveManager) GetTopLevelScores(n int) []LevelScoreEntry {
	data := sm.saves[sm.currentSlot]
	if data == nil || data.LevelScores == nil {
		return []LevelScoreEntry{}
	}

	// Convert map to slice
	entries := make([]LevelScoreEntry, 0, len(data.LevelScores))
	for levelID, score := range data.LevelScores {
		entries = append(entries, LevelScoreEntry{
			LevelID: levelID,
			Score:   score,
		})
	}

	// Sort by score descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	// Return top N
	if n > len(entries) {
		n = len(entries)
	}
	return entries[:n]
}

// WorldScoreEntry represents a world score entry for leaderboards
type WorldScoreEntry struct {
	WorldID int
	Score   int
}

// GetTopWorldScores returns the top N world scores
func (sm *SaveManager) GetTopWorldScores(n int) []WorldScoreEntry {
	data := sm.saves[sm.currentSlot]
	if data == nil || data.WorldScores == nil {
		return []WorldScoreEntry{}
	}

	// Convert map to slice
	entries := make([]WorldScoreEntry, 0, len(data.WorldScores))
	for worldID, score := range data.WorldScores {
		entries = append(entries, WorldScoreEntry{
			WorldID: worldID,
			Score:   score,
		})
	}

	// Sort by score descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	// Return top N
	if n > len(entries) {
		n = len(entries)
	}
	return entries[:n]
}

// GetAllLevelScores returns all level scores for the current save
func (sm *SaveManager) GetAllLevelScores() map[string]int {
	data := sm.saves[sm.currentSlot]
	if data == nil || data.LevelScores == nil {
		return make(map[string]int)
	}

	// Return a copy to prevent external modification
	scores := make(map[string]int, len(data.LevelScores))
	for k, v := range data.LevelScores {
		scores[k] = v
	}
	return scores
}

// GetAllWorldScores returns all world scores for the current save
func (sm *SaveManager) GetAllWorldScores() map[int]int {
	data := sm.saves[sm.currentSlot]
	if data == nil || data.WorldScores == nil {
		return make(map[int]int)
	}

	// Return a copy to prevent external modification
	scores := make(map[int]int, len(data.WorldScores))
	for k, v := range data.WorldScores {
		scores[k] = v
	}
	return scores
}

// ClearLevelScore clears the high score for a specific level
func (sm *SaveManager) ClearLevelScore(levelID string) error {
	data := sm.saves[sm.currentSlot]
	if data == nil {
		return fmt.Errorf("no save data in current slot %d", sm.currentSlot)
	}

	if data.LevelScores != nil {
		delete(data.LevelScores, levelID)
		return sm.SaveGame(sm.currentSlot, data)
	}

	return nil
}

// ClearAllScores clears all high scores for the current save
func (sm *SaveManager) ClearAllScores() error {
	data := sm.saves[sm.currentSlot]
	if data == nil {
		return fmt.Errorf("no save data in current slot %d", sm.currentSlot)
	}

	data.LevelScores = make(map[string]int)
	data.WorldScores = make(map[int]int)
	data.TotalScore = 0

	return sm.SaveGame(sm.currentSlot, data)
}
