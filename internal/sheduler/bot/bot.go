package bot

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"log"
	"strconv"
	"sync"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot            	*tgbotapi.BotAPI
	UserRepository 	repositories.UserRepository
	EventRepository	repositories.EventRepository

}

func CreateBot(stop chan struct{}, userRepository repositories.UserRepository, eventRepository repositories.EventRepository, token string) (*Bot,error) {
	var err error
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil,err
	}
	Bot := Bot{bot: bot, UserRepository: userRepository,EventRepository: eventRepository}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Bot.listenForUpdates(stop)
	}()
	return &Bot,nil
}

func (b *Bot) SendMsg(event entities.Event, msg string) error {
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()
	members,err:=b.EventRepository.FetchMembers(ctx,event.Id.String()) 
	if err!=nil{
		return err
	}
	for _, id := range members{
		user, err := b.UserRepository.FindById(context.Background(), id)
		if err != nil {
			return err
		}
		if user.ChatId != "" {
			chatID, _ := strconv.ParseInt(user.ChatId, 10, 64)
			message := tgbotapi.NewMessage(chatID, msg)
			if _, err := b.bot.Send(message); err != nil {
				log.Printf("failed to send message to user %s: %v", user.Telegram, err)
			}
		}
	}
	return nil
}

func (b *Bot) listenForUpdates(stop chan struct{}) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := b.bot.GetUpdatesChan(updateConfig)
	for {
		select {
		case <-stop:
			log.Println("stopping bot")
			return
		case update := <-updates:
			if update.Message != nil {
				b.handleMessage(update)
			}
		}
	}

}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	username := update.Message.From.UserName
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()
	user, err := b.UserRepository.FindBy(ctx,"telegram",username)
	if err != nil {
		log.Printf("user not found: %v", err)
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅ Да, хочу"),
			tgbotapi.NewKeyboardButton("❌ Нет, не хочу"),
		),
	)

	switch text {
	case "✅ Да, хочу":
		if user.ChatId == "" {
			if err := b.storeChatID(user, chatID); err != nil {
				log.Printf("failed to store chatId: %v", err)
			}
			msg := tgbotapi.NewMessage(chatID, "Теперь вы будете получать уведомления от crap о начале ивентов!")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		} else {
			msg := tgbotapi.NewMessage(chatID, "Вы уже подписаны на уведомления.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}

	case "❌ Нет, не хочу":
		if user.ChatId != "" {
			if err := b.removeChatID(user); err != nil {
				log.Printf("failed to remove chatId: %v", err)
			}
			msg := tgbotapi.NewMessage(chatID, "Вы отписались от уведомлений.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		} else {
			msg := tgbotapi.NewMessage(chatID, "Вы не подписаны на уведомления.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}

	default:
		msg := tgbotapi.NewMessage(chatID, "Хотите ли вы получать уведомления о начале ивентов, к которым вы присоединились?")
		msg.ReplyMarkup = keyboard
		if _, err := b.bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}

func (b *Bot) storeChatID(user *entities.User, chatID int64) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	user.ChatId = strconv.Itoa(int(chatID))
	if err := b.UserRepository.Save(context.Background(), *user); err != nil {
		return err
	}
	return nil
}

func (b *Bot) removeChatID(user *entities.User) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	user.ChatId = ""
	if err := b.UserRepository.Save(context.Background(), *user); err != nil {
		return err
	}
	return nil
}