package openai_compatible

import (
	"context"

	"github.com/danielmiessler/fabric/internal/plugins/ai/openai"
)

// AllowedTranscriptionModels lists the transcription models supported by
// OpenAI-compatible providers such as SiliconCloud. SenseVoiceSmall is a free
// speech-to-text model; pass --transcribe-language (e.g. zh) to suppress
// language markers in its output.
var AllowedTranscriptionModels = []string{
	"FunAudioLLM/SenseVoiceSmall",
}

// TranscribeFile transcribes the given audio file using the provider's
// OpenAI-compatible audio transcription endpoint. It shadows the promoted
// method from the embedded openai.Client so that the model allowlist and
// language handling apply to this provider class.
func (c *Client) TranscribeFile(ctx context.Context, filePath, model string, split bool, language string) (string, error) {
	return openai.TranscribeAudioFile(ctx, c.ApiClient, filePath, model, split, language, AllowedTranscriptionModels)
}
