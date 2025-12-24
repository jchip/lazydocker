package i18n

import (
	"strings"

	"github.com/imdario/mergo"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
)

// Localizer will translate a message into the user's language
type Localizer struct {
	Log *logrus.Entry
	S   TranslationSet
}

func NewTranslationSetFromConfig(log *logrus.Entry, configLanguage string, nerdFontsVersion string) (*TranslationSet, error) {
	if configLanguage == "auto" {
		language := detectLanguage(jibber_jabber.DetectLanguage)
		return NewTranslationSet(log, language, nerdFontsVersion), nil
	}

	for key := range GetTranslationSets() {
		if key == configLanguage {
			return NewTranslationSet(log, configLanguage, nerdFontsVersion), nil
		}
	}

	return NewTranslationSet(log, "en", nerdFontsVersion), errors.New("Language not found: " + configLanguage)
}

func NewTranslationSet(log *logrus.Entry, language string, nerdFontsVersion string) *TranslationSet {
	log.Info("language: " + language)

	baseSet := englishSet()

	for languageCode, translationSet := range GetTranslationSets() {
		if strings.HasPrefix(language, languageCode) {
			_ = mergo.Merge(&baseSet, translationSet, mergo.WithOverride)
		}
	}

	if nerdFontsVersion != "" {
		// Icons from: https://www.nerdfonts.com/cheat-sheet
		baseSet.ServicesTitle = " Services"     // nf-fa-server
		baseSet.ContainersTitle = " Containers" // nf-fa-cube
		baseSet.StandaloneContainersTitle = " Standalone Containers"
		baseSet.ImagesTitle = " Images"     // nf-linux-docker
		baseSet.VolumesTitle = " Volumes"   // nf-fa-database
		baseSet.NetworksTitle = " Networks" // nf-fa-sitemap
		baseSet.ConfigTitle = " Config"     // nf-oct-gear
		baseSet.ProjectTitle = " Project"   // nf-fa-folder_open
		baseSet.ContainerConfigTitle = " Container Config"
	}

	return &baseSet
}

// GetTranslationSets gets all the translation sets, keyed by language code
func GetTranslationSets() map[string]TranslationSet {
	return map[string]TranslationSet{
		"pl": polishSet(),
		"nl": dutchSet(),
		"de": germanSet(),
		"tr": turkishSet(),
		"en": englishSet(),
		"fr": frenchSet(),
		"zh": chineseSet(),
		"es": spanishSet(),
		"pt": portugueseSet(),
	}
}

// detectLanguage extracts user language from environment
func detectLanguage(langDetector func() (string, error)) string {
	if userLang, err := langDetector(); err == nil {
		return userLang
	}

	return "C"
}
