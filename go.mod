module kube-ai

go 1.17

require (
	github.com/joho/godotenv v1.5.1 // .env dosyasını okumak için
	github.com/sashabaranov/go-openai v1.38.1 // OpenAI API istemcisi
	github.com/spf13/cobra v1.9.1 // CLI komut yapısı
)

require (
	// aşağıdakiler cobra veya diğer paketler tarafından dolaylı olarak kullanılıyor
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
)
