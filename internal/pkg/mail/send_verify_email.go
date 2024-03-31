package mail

import (
	`context`
	`fmt`
	`log/slog`

	`github.com/dugtriol/BarterApp/internal/config`
)

func SendVerifyEmail(ctx context.Context, log *slog.Logger, name, emailPath string) error {
	cfg := config.MustLoad()
	// TODO: отдельная таблица в БД для сохранения отправленных emails
	subject := "Welcome to Barter App"
	// TODO: replace this URL with an environment variable that points to a front-end page
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=&secret_code=")
	content := fmt.Sprintf(
		`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, name, verifyUrl,
	)
	to := []string{emailPath}

	sender := NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)
	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}
	return nil
}
