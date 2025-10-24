package mail

import (
	"testing"

	"github.com/JeongWoo-Seo/simpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "test email"
	content := `
	<h1>test</h1>
	<p> test test </p>`

	to := []string{"sjwoo6500@gmail.com"}

	err = sender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}

func BenchmarkSendEmailWithGmail(b *testing.B) {
	config, err := util.LoadConfig("..")
	require.NoError(b, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Benchmark Email"
	content := `<h1>Benchmark Test</h1><p>Measuring email sending latency.</p>`
	to := []string{"sjwoo6500@gmail.com"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = sender.SendEmail(subject, content, to, nil, nil, nil)

		if err != nil {
			b.Fatal(err)
		}
	}
}
