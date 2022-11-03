package filechecker

import (
	"mime/multipart"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

const (
	TypeAPPLICATION = "Application"
	TypeARCHIVE     = "Archive"
	TypeAUDIO       = "Audio"
	TypeDOCUMENTS   = "Documents"
	TypeFONT        = "Font"
	TypeIMAGE       = "Image"
	TypeVIDEO       = "Video"

	ExtAppDEX  = "dex"
	ExtAppDEY  = "dey"
	ExtAppWASM = "wasm"

	ExtArchive7Z     = "7z"
	ExtArchiveZ      = "Z"
	ExtArchiveAR     = "ar"
	ExtArchiveBZ2    = "bz2"
	ExtArchiveCAB    = "cab"
	ExtArchiveCRX    = "crx"
	ExtArchiveDCM    = "dcm"
	ExtArchiveDEB    = "deb"
	ExtArchiveELF    = "elf"
	ExtArchiveEOT    = "eot"
	ExtArchiveEPUB   = "epub"
	ExtArchiveEXE    = "exe"
	ExtArchiveGZ     = "gz"
	ExtArchiveISO    = "iso"
	ExtArchiveLZ     = "lz"
	ExtArchiveNES    = "nes"
	ExtArchivePDF    = "pdf"
	ExtArchivePS     = "ps"
	ExtArchiveRAR    = "rar"
	ExtArchiveRPM    = "rpm"
	ExtArchiveRTF    = "rtf"
	ExtArchiveSQLITE = "sqlite"
	ExtArchiveSWF    = "swf"
	ExtArchiveTAR    = "tar"
	ExtArchiveXZ     = "xz"
	ExtArchiveZIP    = "zip"
	ExtArchiveZSTD   = "zstd"

	ExtAudioAAC  = "aac"
	ExtAudioAIFF = "aiff"
	ExtAudioAMR  = "amr"
	ExtAudioFLAC = "flac"
	ExtAudioM4A  = "m4a"
	ExtAudioMID  = "mid"
	ExtAudioMP3  = "mp3"
	ExtAudioOGG  = "ogg"
	ExtAudioWAV  = "wav"

	ExtDocDOC  = "doc"
	ExtDocDOCX = "docx"
	ExtDocPPT  = "ppt"
	ExtDocPPTX = "pptx"
	ExtDocXLS  = "xls"
	ExtDocXLSX = "xlsx"

	ExtFontOTF   = "otf"
	ExtFontTTF   = "ttf"
	ExtFontWOFF  = "woff"
	ExtFontWOFF2 = "woff2"

	ExtImgBMP  = "bmp"
	ExtImgCR2  = "cr2"
	ExtImgDWG  = "dwg"
	ExtImgGIF  = "gif"
	ExtImgHEIF = "heif"
	ExtImgICO  = "ico"
	ExtImgJPG  = "jpg"
	ExtImgJXR  = "jxr"
	ExtImgPNG  = "png"
	ExtImgPSD  = "psd"
	ExtImgTIF  = "tif"
	ExtImgWEBP = "webp"

	ExtVideo3GP  = "3gp"
	ExtVideoAVI  = "avi"
	ExtVideoFLV  = "flv"
	ExtVideoM4V  = "m4v"
	ExtVideoMKV  = "mkv"
	ExtVideoMOV  = "mov"
	ExtVideoMP4  = "mp4"
	ExtVideoMPG  = "mpg"
	ExtVideoWEBM = "webm"
	ExtVideoWMV  = "wmv"
)

type FileChecker struct {
	file *multipart.FileHeader

	// list of all authorised types, computed. authorisedTypes[typ] = true|false
	authorisedTypes map[string]bool

	// list of all authorised extensions (defaults + those that have been
	// user-authorised/unauthorised). authorisedExtensions[ext] = true|false
	authorisedExtensions map[string]bool

	// a dictionary of extensions with their corresponding file types.
	// _dictionary[ext] = typ
	_dictionary map[string]string
}

var (
	availableExtensions = map[string]map[string]bool{
		TypeAPPLICATION: {
			ExtAppDEX:  false,
			ExtAppDEY:  false,
			ExtAppWASM: false,
		},

		TypeARCHIVE: {
			ExtArchive7Z:     false,
			ExtArchiveZ:      false,
			ExtArchiveAR:     false,
			ExtArchiveBZ2:    false,
			ExtArchiveCAB:    false,
			ExtArchiveCRX:    false,
			ExtArchiveDCM:    false,
			ExtArchiveDEB:    false,
			ExtArchiveELF:    false,
			ExtArchiveEOT:    false,
			ExtArchiveEPUB:   false,
			ExtArchiveEXE:    false,
			ExtArchiveGZ:     false,
			ExtArchiveISO:    false,
			ExtArchiveLZ:     false,
			ExtArchiveNES:    false,
			ExtArchivePDF:    true, // allowed by default.
			ExtArchivePS:     false,
			ExtArchiveRAR:    false,
			ExtArchiveRPM:    false,
			ExtArchiveRTF:    false,
			ExtArchiveSQLITE: false,
			ExtArchiveSWF:    false,
			ExtArchiveTAR:    false,
			ExtArchiveXZ:     false,
			ExtArchiveZIP:    false,
			ExtArchiveZSTD:   false,
		},

		TypeAUDIO: {
			ExtAudioAAC:  false,
			ExtAudioAIFF: false,
			ExtAudioAMR:  false,
			ExtAudioFLAC: false,
			ExtAudioM4A:  false,
			ExtAudioMID:  false,
			ExtAudioMP3:  false,
			ExtAudioOGG:  false,
			ExtAudioWAV:  false,
		},

		TypeDOCUMENTS: {
			ExtDocDOC:  false,
			ExtDocDOCX: false,
			ExtDocPPT:  false,
			ExtDocPPTX: false,
			ExtDocXLS:  false,
			ExtDocXLSX: false,
		},

		TypeFONT: {
			ExtFontOTF:   false,
			ExtFontTTF:   false,
			ExtFontWOFF:  false,
			ExtFontWOFF2: false,
		},

		TypeIMAGE: {
			ExtImgBMP:  false,
			ExtImgCR2:  false,
			ExtImgDWG:  false,
			ExtImgGIF:  false,
			ExtImgHEIF: false,
			ExtImgICO:  false,
			ExtImgJPG:  true, // allowed by default.
			ExtImgJXR:  false,
			ExtImgPNG:  true, // allowed by default.
			ExtImgPSD:  false,
			ExtImgTIF:  false,
			ExtImgWEBP: false,
		},

		TypeVIDEO: {
			ExtVideo3GP:  false,
			ExtVideoAVI:  false,
			ExtVideoFLV:  false,
			ExtVideoM4V:  false,
			ExtVideoMKV:  false,
			ExtVideoMOV:  false,
			ExtVideoMP4:  false,
			ExtVideoMPG:  false,
			ExtVideoWEBM: false,
			ExtVideoWMV:  false,
		},
	}
)

// GetFileChecker returns an instance of FileChecker.
func GetFileChecker(file *multipart.FileHeader) *FileChecker {
	var (
		authorisedTypes      = make(map[string]bool, 0)
		authorisedExtensions = make(map[string]bool, 0)
		_dictionary          = make(map[string]string, 0)
	)

	// default authorised extensions
	for typ, extensions := range availableExtensions {
		for ext, authorised := range extensions {
			if authorised {
				authorisedTypes[typ] = true
				authorisedExtensions[ext] = true
			}
			_dictionary[ext] = typ
		}
	}

	return &FileChecker{
		file:                 file,
		authorisedTypes:      authorisedTypes,
		authorisedExtensions: authorisedExtensions,
		_dictionary:          _dictionary,
	}
}

// SetFile sets the file to be checked.
func (fc *FileChecker) SetFile(file *multipart.FileHeader) {
	if file != nil {
		fc.file = file
	}
}

// SetExtensions sets authorised extensions.
func (fc *FileChecker) SetExtensions(extensions []string) {
	for _, ext := range extensions {
		if typ, found := fc._dictionary[ext]; found {
			fc.authorisedTypes[typ] = true
			fc.authorisedExtensions[ext] = true
		}
	}
}

// UnsetExtensions unsets types that were authorised.
func (fc *FileChecker) UnsetExtensions(extensions []string) {
	// un-authorise the extensions we've been requested
	for _, ext := range extensions {
		fc.authorisedExtensions[ext] = false
	}

	// we un-authorised all types (to re-authorise later below)
	for typ := range fc.authorisedTypes {
		fc.authorisedTypes[typ] = false
	}

	// for all extensions, find those authorised, and hence authorise the
	// corresponding type(s)
	for ext, authorised := range fc.authorisedExtensions {
		if authorised {
			typ, _ := fc._dictionary[ext]
			fc.authorisedTypes[typ] = true
		}
	}
}

// IsAuthorised tells us whether the file is authorised (type and extension).
func (fc *FileChecker) IsAuthorised() bool {
	var (
		err  error
		file multipart.File
		kind types.Type

		// file header, first 261 bytes (to be read further below)
		// see https://www.garykessler.net/library/file_sigs.html
		header = make([]byte, 261)
	)

	// file was not provided or wrongly provided
	if fc.file == nil {
		return false
	}

	// cannot open
	if file, err = fc.file.Open(); err != nil {
		return false
	}

	// cannot read header
	if _, err = file.Read(header); err != nil {
		return false
	}

	// cannot match header
	if kind, err = filetype.Match(header); err != nil {
		return false
	}

	// verify authorised types
	if authorised := fc.isTypeAuthorised(header); !authorised {
		return false
	}

	// extension not among those available or available extension is not
	// authorised (set to false)
	if authorised, found := fc.authorisedExtensions[kind.Extension]; !authorised || !found {
		return false
	}

	return true
}

// isTypeAuthorised is a private method. Checks if type of file is authorised.
func (fc *FileChecker) isTypeAuthorised(header []byte) bool {
	if _, authorised := fc.authorisedTypes[TypeAPPLICATION]; authorised && filetype.IsApplication(header) {
		return true
	}

	if _, authorised := fc.authorisedTypes[TypeARCHIVE]; authorised && filetype.IsArchive(header) {
		return true
	}

	if _, authorised := fc.authorisedTypes[TypeAUDIO]; authorised && filetype.IsAudio(header) {
		return true
	}

	if _, authorised := fc.authorisedTypes[TypeDOCUMENTS]; authorised && filetype.IsDocument(header) {
		return true
	}

	if _, authorised := fc.authorisedTypes[TypeFONT]; authorised && filetype.IsFont(header) {
		return true
	}

	if _, authorised := fc.authorisedTypes[TypeIMAGE]; authorised && filetype.IsImage(header) {
		return true
	}

	if _, authorised := fc.authorisedTypes[TypeVIDEO]; authorised && filetype.IsVideo(header) {
		return true
	}

	return false
}
