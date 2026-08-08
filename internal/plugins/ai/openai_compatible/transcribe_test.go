package openai_compatible

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/danielmiessler/fabric/internal/i18n"
	"github.com/danielmiessler/fabric/internal/plugins/ai/openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranscribeFile_ModelNotSupportedIsLocalized(t *testing.T) {
	_, err := i18n.Init("en")
	require.NoError(t, err)

	audioFile, err := os.CreateTemp("", "transcribe-valid-*.mp3")
	require.NoError(t, err)
	require.NoError(t, audioFile.Close())
	t.Cleanup(func() { _ = os.Remove(audioFile.Name()) })

	// whisper-1 is in the OpenAI allowlist but not the OpenAI-compatible one.
	_, err = openai.TranscribeAudioFile(context.Background(), nil, audioFile.Name(), "whisper-1", false, "", AllowedTranscriptionModels)
	require.Error(t, err)
	assert.Equal(t,
		fmt.Sprintf(i18n.T("openai_audio_model_not_supported_for_transcription"), "whisper-1"),
		err.Error(),
	)
}

func TestTranscribeFile_SenseVoiceSmallIsAllowed(t *testing.T) {
	_, err := i18n.Init("en")
	require.NoError(t, err)

	// SenseVoiceSmall passes the allowlist; the error must come from the
	// unsupported extension check, which runs before any API call.
	unsupportedFile, err := os.CreateTemp("", "transcribe-sensevoice-*.txt")
	require.NoError(t, err)
	require.NoError(t, unsupportedFile.Close())
	t.Cleanup(func() { _ = os.Remove(unsupportedFile.Name()) })

	_, err = openai.TranscribeAudioFile(context.Background(), nil, unsupportedFile.Name(), AllowedTranscriptionModels[0], false, "zh", AllowedTranscriptionModels)
	require.Error(t, err)
	assert.Equal(t,
		fmt.Sprintf(i18n.T("openai_audio_unsupported_audio_format"), filepath.Ext(unsupportedFile.Name())),
		err.Error(),
	)
}

func TestTranscribeFile_SharedFunctionHonorsNilAllowedModels(t *testing.T) {
	_, err := i18n.Init("en")
	require.NoError(t, err)

	// nil allowlist accepts any model; the oversized-file check must fire
	// before any API call.
	largeFile, err := os.CreateTemp("", "transcribe-large-*.mp3")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Remove(largeFile.Name()) })
	require.NoError(t, largeFile.Truncate(openai.MaxAudioFileSize+1))
	require.NoError(t, largeFile.Close())

	_, err = openai.TranscribeAudioFile(context.Background(), nil, largeFile.Name(), "any-model", false, "", nil)
	require.Error(t, err)
	assert.Equal(t,
		fmt.Sprintf(i18n.T("openai_audio_file_exceeds_limit_enable_split"), largeFile.Name()),
		err.Error(),
	)
}
