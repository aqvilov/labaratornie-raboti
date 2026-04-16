package main

//пишем свой спотифай
import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// струткура трека

type Track struct {
	Name     string  `json:"name"`
	Duration int     `json:"duration"`
	Genre    string  `json:"genre"`
	Rating   float64 `json:"rating"`
}

func (t *Track) String() string {
	minutes := t.Duration / 60
	seconds := t.Duration % 60
	return fmt.Sprintf("%s | %02d:%02d | %s | %.1f", t.Name, minutes, seconds, t.Genre, t.Rating)
}

type RepeatMode int

const (
	RepeatNone RepeatMode = iota
	RepeatOne
	RepeatAll
)

func (r RepeatMode) String() string {
	switch r {
	case RepeatNone:
		return "none"
	case RepeatOne:
		return "one"
	case RepeatAll:
		return "all"
	default:
		return "unknown"
	}
}

// структура плейлиста
type Playlist struct {
	Name         string
	tracks       []*Track
	currentIndex int
	repeatMode   RepeatMode
	rng          *rand.Rand //генератор чисел для перемешниваняи
}

// типо конструктора из питона
func NewPlaylist(name string) *Playlist {
	return &Playlist{
		Name:         name,
		tracks:       make([]*Track, 0),
		currentIndex: 0,
		repeatMode:   RepeatNone,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (p *Playlist) AddTrack(track *Track) {
	p.tracks = append(p.tracks, track)
	fmt.Printf("Трек '%s' добавлен в плейлист '%s'\n", track.Name, p.Name)
}

func (p *Playlist) DeleteTrack(trackName string) bool {
	for i, track := range p.tracks {
		if track.Name == trackName {
			p.tracks = append(p.tracks[:i], p.tracks[i+1:]...)
			if p.currentIndex >= len(p.tracks) && len(p.tracks) > 0 {
				p.currentIndex = len(p.tracks) - 1
			} else if len(p.tracks) == 0 {
				p.currentIndex = 0
			}
			fmt.Printf("трек '%s' удален из плейлиста\n", trackName)
			return true
		}
	}
	fmt.Printf("Трек '%s' не найден в плейлисте\n", trackName)
	return false
}

// перемешка треков
func (p *Playlist) Shuffle() {
	if len(p.tracks) == 0 {
		fmt.Println("Плейлист пуст, нечего перемешивать")
		return
	}

	p.rng.Shuffle(len(p.tracks), func(i, j int) {
		p.tracks[i], p.tracks[j] = p.tracks[j], p.tracks[i]
	})
	p.currentIndex = 0
	fmt.Println("Плейлист перемешан")
}

func (p *Playlist) SetRepeatMode(mode RepeatMode) {
	p.repeatMode = mode
	fmt.Printf("Режим повтора установлен на: %s\n", mode)
}

func (p *Playlist) CurrentTrack() *Track {
	if len(p.tracks) == 0 {
		return nil
	}
	if p.currentIndex >= 0 && p.currentIndex < len(p.tracks) {
		return p.tracks[p.currentIndex]
	}
	return nil
}

func (p *Playlist) NextTrack() *Track {
	if len(p.tracks) == 0 {
		fmt.Println("Плейлист пуст")
		return nil
	}

	if p.repeatMode == RepeatOne {
		return p.CurrentTrack()
	}

	nextIndex := p.currentIndex + 1
	if nextIndex < len(p.tracks) {
		p.currentIndex = nextIndex
	} else if p.repeatMode == RepeatAll {
		p.currentIndex = 0
		fmt.Println("Плейлист зациклен, воспроизведение с начала")
	} else {
		p.currentIndex = 0
		fmt.Println("Достигнут конец плейлиста")
	}

	return p.CurrentTrack()
}

func (p *Playlist) PrevTrack() *Track {
	if len(p.tracks) == 0 {
		fmt.Println("Плейлист пуст")
		return nil
	}

	if p.repeatMode == RepeatOne {
		return p.CurrentTrack()
	}

	prevIndex := p.currentIndex - 1
	if prevIndex >= 0 {
		p.currentIndex = prevIndex
	} else if p.repeatMode == RepeatAll {
		p.currentIndex = len(p.tracks) - 1
		fmt.Println("Плейлист зациклен, воспроизведение с конца")
	} else {
		p.currentIndex = 0
		fmt.Println("Достигнуто начало плейлиста")
	}

	return p.CurrentTrack()
}

// отображение вскх треков из плейлиста
func (p *Playlist) Display() {
	fmt.Printf("\n   Плейлист: %s   \n", p.Name)
	fmt.Printf("Режим повтора: %s\n", p.repeatMode)
	fmt.Println(strings.Repeat("-", 80))

	headers := []string{"№", "Название", "Длительность", "Жанр", "Рейтинг"}
	fmt.Printf("%-4s | %-30s | %-10s | %-15s | %-6s\n", headers[0], headers[1], headers[2], headers[3], headers[4])
	fmt.Println(strings.Repeat("-", 80))

	for i, track := range p.tracks {
		current := ""
		if i == p.currentIndex {
			current = "Play"
		} else {
			current = " "
		}

		minutes := track.Duration / 60
		seconds := track.Duration % 60
		durationStr := fmt.Sprintf("%02d:%02d", minutes, seconds)

		fmt.Printf("%s%-3d | %-30s | %-10s | %-15s | %-6.1f\n",
			current, i+1, track.Name, durationStr, track.Genre, track.Rating)
	}
	fmt.Println(strings.Repeat("-", 80))

	if current := p.CurrentTrack(); current != nil {
		fmt.Printf("Сейчас играет: %s\n", current.Name)
	}
}

func (p *Playlist) FindTracksInTimeRange(n int, minDuration, maxDuration int) []string {
	result := make([]string, 0)
	for _, track := range p.tracks {
		if track.Duration >= minDuration && track.Duration <= maxDuration {
			result = append(result, track.Name)
			if len(result) == n {
				break
			}
		}
	}
	return result
}

func (p *Playlist) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, track := range p.tracks {
		trackData, err := json.Marshal(track)
		if err != nil {
			return fmt.Errorf("ошибка сериализации трека: %v", err)
		}
		encoded := base64.StdEncoding.EncodeToString(trackData)

		_, err = writer.WriteString(encoded + "\n")
		if err != nil {
			return fmt.Errorf("ошибка записи в файл: %v", err)
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("ошибка сохранения файла: %v", err)
	}

	fmt.Printf("Плейлист успешно сохранен в файл: %s\n", filename)
	return nil
}

func (p *Playlist) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	newTracks := make([]*Track, 0)

	for scanner.Scan() {
		encoded := scanner.Text()
		if encoded == "" {
			continue
		}

		//декодируем из басе64
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return fmt.Errorf("ошибка декодирования Base64: %v", err)
		}

		var track Track
		err = json.Unmarshal(decoded, &track)
		if err != nil {
			return fmt.Errorf("ошибка десериализации трека: %v", err)
		}

		newTracks = append(newTracks, &track)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	p.tracks = newTracks
	p.currentIndex = 0
	fmt.Printf("Плейлист успешно загружен из файла: %s (%d треков)\n", filename, len(p.tracks))
	return nil
}

func main() {
	playlist := NewPlaylist("Лучший плейлист")

	tracks := []*Track{
		{Name: "Come as you are", Duration: 354, Genre: "Rock", Rating: 8.3},
		{Name: "Самый лучший эмо панк", Duration: 282, Genre: "Rock", Rating: 9.9},
		{Name: "Кайен", Duration: 183, Genre: "Pop", Rating: 9.1},
		{Name: "Улица сталеваров", Duration: 391, Genre: "Pop", Rating: 7.6},
		{Name: "Heart-ShapedBox", Duration: 294, Genre: "Grunge", Rating: 9.2},
		{Name: "Smells Like Teen Spirit", Duration: 301, Genre: "Grunge", Rating: 10.0},
	}

	for _, track := range tracks {
		playlist.AddTrack(track)
	}
	// отображение на экране
	playlist.Display()

	fmt.Println("\n--- Навигация ---")
	fmt.Printf("Текущий трек: %s\n", playlist.CurrentTrack().Name)

	playlist.NextTrack()
	fmt.Printf("Следующий трек: %s\n", playlist.CurrentTrack().Name)

	playlist.PrevTrack()
	fmt.Printf("Предыдущий трек: %s\n", playlist.CurrentTrack().Name)

	fmt.Println("\n--- Перемешивание ---")
	// перемешиваем + отображаем
	playlist.Shuffle()
	playlist.Display()

	fmt.Println("\n--- Режимы повтора ---")
	playlist.SetRepeatMode(RepeatOne)
	fmt.Printf("Повтор одного трека: %s\n", playlist.CurrentTrack().Name)
	playlist.NextTrack()
	fmt.Printf("После next с repeat one: %s\n", playlist.CurrentTrack().Name)

	playlist.SetRepeatMode(RepeatAll)
	playlist.NextTrack()

	fmt.Println("\n--- Поиск треков ---")
	found := playlist.FindTracksInTimeRange(3, 180, 360)
	fmt.Printf("Найдено треков длительностью от 3 до 6 минут: %v\n", found)

	fmt.Println("\n--- Удаление трека ---")
	playlist.DeleteTrack("Улица сталеваров")
	playlist.Display()

	fmt.Println("\n--- Сохранение и загрузка ---")
	err := playlist.SaveToFile("playlist.txt")
	if err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
	}

	newPlaylist := NewPlaylist("Загруженный плейлист")
	err = newPlaylist.LoadFromFile("playlist.txt")
	if err != nil {
		fmt.Printf("ошибка загрузки: %v\n", err)
	} else {
		newPlaylist.Display()
	}
}
