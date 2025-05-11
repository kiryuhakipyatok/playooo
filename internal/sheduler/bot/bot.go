package bot

import (
	"context"
	"crap/internal/domain/entities"
	"crap/internal/domain/repositories"
	"strconv"
	"sync"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	Logger          *logrus.Logger
	UserRepository  repositories.UserRepository
	EventRepository repositories.EventRepository
}

func CreateBot(stop chan struct{}, l *logrus.Logger, userRepository repositories.UserRepository, eventRepository repositories.EventRepository, token string) (*Bot, error) {
	var err error
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	Bot := Bot{bot: bot, UserRepository: userRepository, EventRepository: eventRepository, Logger: l}
	return &Bot, err
}

func (b *Bot) SendMsg(event entities.Event, msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	members, err := b.EventRepository.FetchMembers(ctx, event.Id.String())
	if err != nil {
		return err
	}
	for _, id := range members {
		user, err := b.UserRepository.FindById(context.Background(), id)
		if err != nil {
			return err
		}
		if user.ChatId != "" {
			chatID, _ := strconv.ParseInt(user.ChatId, 10, 64)
			message := tgbotapi.NewMessage(chatID, msg)
			if _, err := b.bot.Send(message); err != nil {
				b.Logger.Infof("failed to send message to user %s: %v", user.Telegram, err)
			}
		}
	}
	return nil
}

func (b *Bot) ListenForUpdates(stop chan struct{}) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := b.bot.GetUpdatesChan(updateConfig)
	for {
		select {
		case <-stop:
			b.Logger.Info("stopping bot")
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user, err := b.UserRepository.FindBy(ctx, "telegram", username)
	if err != nil {
		b.Logger.WithError(err).Info("user not found")
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
				b.Logger.WithError(err).Info("failed to store chatID")
			}
			msg := tgbotapi.NewMessage(chatID, "Теперь вы будете получать уведомления от crap о начале ивентов!")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				b.Logger.WithError(err).Info("failed to send msg")
			}
		} else {
			msg := tgbotapi.NewMessage(chatID, "Вы уже подписаны на уведомления.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				b.Logger.WithError(err).Info("failed to send msg")
			}
		}

	case "❌ Нет, не хочу":
		if user.ChatId != "" {
			if err := b.removeChatID(user); err != nil {
				b.Logger.WithError(err).Info("failed to remove chatID")
			}
			msg := tgbotapi.NewMessage(chatID, "Вы отписались от уведомлений.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				b.Logger.WithError(err).Info("failed to send msg")
			}
		} else {
			msg := tgbotapi.NewMessage(chatID, "Вы не подписаны на уведомления.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				b.Logger.WithError(err).Info("failed to send msg")
			}
		}

	default:
		msg := tgbotapi.NewMessage(chatID, "Хотите ли вы получать уведомления о начале ивентов, к которым вы присоединились?")
		msg.ReplyMarkup = keyboard
		if _, err := b.bot.Send(msg); err != nil {
			b.Logger.WithError(err).Info("failed to send msg")
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
