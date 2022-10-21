package filechecker

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

var (
	defaultAuthorisedTypes = map[string]bool{
		TypeIMAGE:   true,
		TypeARCHIVE: true,
	}

	defaultAuthorisedExtensions = map[string]bool{
		ExtArchivePDF: true,
		ExtImgJPG:     true,
		ExtImgPNG:     true,
	}

	defaultDictionary = map[string]string{
		ExtAppDEX:  TypeAPPLICATION,
		ExtAppDEY:  TypeAPPLICATION,
		ExtAppWASM: TypeAPPLICATION,

		ExtArchive7Z:     TypeARCHIVE,
		ExtArchiveZ:      TypeARCHIVE,
		ExtArchiveAR:     TypeARCHIVE,
		ExtArchiveBZ2:    TypeARCHIVE,
		ExtArchiveCAB:    TypeARCHIVE,
		ExtArchiveCRX:    TypeARCHIVE,
		ExtArchiveDCM:    TypeARCHIVE,
		ExtArchiveDEB:    TypeARCHIVE,
		ExtArchiveELF:    TypeARCHIVE,
		ExtArchiveEOT:    TypeARCHIVE,
		ExtArchiveEPUB:   TypeARCHIVE,
		ExtArchiveEXE:    TypeARCHIVE,
		ExtArchiveGZ:     TypeARCHIVE,
		ExtArchiveISO:    TypeARCHIVE,
		ExtArchiveLZ:     TypeARCHIVE,
		ExtArchiveNES:    TypeARCHIVE,
		ExtArchivePDF:    TypeARCHIVE,
		ExtArchivePS:     TypeARCHIVE,
		ExtArchiveRAR:    TypeARCHIVE,
		ExtArchiveRPM:    TypeARCHIVE,
		ExtArchiveRTF:    TypeARCHIVE,
		ExtArchiveSQLITE: TypeARCHIVE,
		ExtArchiveSWF:    TypeARCHIVE,
		ExtArchiveTAR:    TypeARCHIVE,
		ExtArchiveXZ:     TypeARCHIVE,
		ExtArchiveZIP:    TypeARCHIVE,
		ExtArchiveZSTD:   TypeARCHIVE,

		ExtAudioAAC:  TypeAUDIO,
		ExtAudioAIFF: TypeAUDIO,
		ExtAudioAMR:  TypeAUDIO,
		ExtAudioFLAC: TypeAUDIO,
		ExtAudioM4A:  TypeAUDIO,
		ExtAudioMID:  TypeAUDIO,
		ExtAudioMP3:  TypeAUDIO,
		ExtAudioOGG:  TypeAUDIO,
		ExtAudioWAV:  TypeAUDIO,

		ExtDocDOC:  TypeDOCUMENTS,
		ExtDocDOCX: TypeDOCUMENTS,
		ExtDocPPT:  TypeDOCUMENTS,
		ExtDocPPTX: TypeDOCUMENTS,
		ExtDocXLS:  TypeDOCUMENTS,
		ExtDocXLSX: TypeDOCUMENTS,

		ExtFontOTF:   TypeFONT,
		ExtFontTTF:   TypeFONT,
		ExtFontWOFF:  TypeFONT,
		ExtFontWOFF2: TypeFONT,

		ExtImgBMP:  TypeIMAGE,
		ExtImgCR2:  TypeIMAGE,
		ExtImgDWG:  TypeIMAGE,
		ExtImgGIF:  TypeIMAGE,
		ExtImgHEIF: TypeIMAGE,
		ExtImgICO:  TypeIMAGE,
		ExtImgJPG:  TypeIMAGE,
		ExtImgJXR:  TypeIMAGE,
		ExtImgPNG:  TypeIMAGE,
		ExtImgPSD:  TypeIMAGE,
		ExtImgTIF:  TypeIMAGE,
		ExtImgWEBP: TypeIMAGE,

		ExtVideo3GP:  TypeVIDEO,
		ExtVideoAVI:  TypeVIDEO,
		ExtVideoFLV:  TypeVIDEO,
		ExtVideoM4V:  TypeVIDEO,
		ExtVideoMKV:  TypeVIDEO,
		ExtVideoMOV:  TypeVIDEO,
		ExtVideoMP4:  TypeVIDEO,
		ExtVideoMPG:  TypeVIDEO,
		ExtVideoWEBM: TypeVIDEO,
		ExtVideoWMV:  TypeVIDEO,
	}

	jpgPath  = "./utilities/assets/nadim.jpg"
	pngPath  = "./utilities/assets/nadim.png"
	pdfPath  = "./utilities/assets/nadim.pdf"
	fakePath = "./utilities/assets/fake.png"
)

func TestGetFileChecker(t *testing.T) {
	type args struct {
		file *multipart.FileHeader
	}

	tests := []struct {
		name string
		args args
		want *FileChecker
	}{
		{
			name: "JPG",
			args: args{file: &multipart.FileHeader{Filename: jpgPath}},
			want: &FileChecker{
				file:                 &multipart.FileHeader{Filename: jpgPath},
				authorisedTypes:      defaultAuthorisedTypes,
				authorisedExtensions: defaultAuthorisedExtensions,
				_dictionary:          defaultDictionary,
			},
		},
		{
			name: "PNG",
			args: args{file: &multipart.FileHeader{Filename: pngPath}},
			want: &FileChecker{
				file:                 &multipart.FileHeader{Filename: pngPath},
				authorisedTypes:      defaultAuthorisedTypes,
				authorisedExtensions: defaultAuthorisedExtensions,
				_dictionary:          defaultDictionary,
			},
		},
		{
			name: "PDF",
			args: args{file: &multipart.FileHeader{Filename: pdfPath}},
			want: &FileChecker{
				file:                 &multipart.FileHeader{Filename: pdfPath},
				authorisedTypes:      defaultAuthorisedTypes,
				authorisedExtensions: defaultAuthorisedExtensions,
				_dictionary:          defaultDictionary,
			},
		},
		{
			name: "FAKE",
			args: args{file: &multipart.FileHeader{Filename: fakePath}},
			want: &FileChecker{
				file:                 &multipart.FileHeader{Filename: fakePath},
				authorisedTypes:      defaultAuthorisedTypes,
				authorisedExtensions: defaultAuthorisedExtensions,
				_dictionary:          defaultDictionary,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileChecker(tt.args.file); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileChecker() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFileChecker_SetFile(t *testing.T) {
	tests := []struct {
		name string
		args *multipart.FileHeader
		want *multipart.FileHeader
	}{
		{
			name: "JPG",
			args: &multipart.FileHeader{Filename: jpgPath},
			want: &multipart.FileHeader{Filename: jpgPath},
		},
		{
			name: "PNG",
			args: &multipart.FileHeader{Filename: pngPath},
			want: &multipart.FileHeader{Filename: pngPath},
		},
		{
			name: "PDF",
			args: &multipart.FileHeader{Filename: pdfPath},
			want: &multipart.FileHeader{Filename: pdfPath},
		},
		{
			name: "FAKE",
			args: &multipart.FileHeader{Filename: fakePath},
			want: &multipart.FileHeader{Filename: fakePath},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := GetFileChecker(nil)
			fc.SetFile(tt.args)
			if !reflect.DeepEqual(fc.file, tt.want) {
				t.Errorf("SetFile() sets %+v, different to %+v", tt.args, tt.want)
			}
		})
	}
}

// nolint:funlen
func TestFileChecker_SetExtensions(t *testing.T) {
	// TEST 1: set the ExtImgWEBP extension.
	// In addition to the default extensions and types set, we expect the
	// ExtImgWEBP extension to be set also (types remain unchanged)
	t.Run("WEBP", func(t *testing.T) {
		var (
			args    = []string{ExtImgWEBP}
			wantExt = getExt([]string{ExtImgWEBP})
			wantTyp = getTyp([]string{})
		)

		fc := GetFileChecker(nil)
		fc.SetExtensions(args)

		if !assert.IsEqual(fc.authorisedExtensions, wantExt) {
			t.Errorf("Set(%+v). Extensions :: Got: %+v. Expected: %+v", args, fc.authorisedExtensions, wantExt)
		}
		if !assert.IsEqual(fc.authorisedTypes, wantTyp) {
			t.Errorf("Set(%+v). Types :: Got: %+v. Expected: %+v", args, fc.authorisedTypes, wantTyp)
		}
	})

	// TEST 2: set the ExtAppDEX extension
	// in addition to the default extensions and types set, we expect the
	// ExtAppDEX extension to be set as well as the TypeAPPLICATION type
	t.Run("DEX", func(t *testing.T) {
		var (
			args    = []string{ExtAppDEX}
			wantExt = getExt([]string{ExtAppDEX})
			wantTyp = getTyp([]string{TypeAPPLICATION})
		)

		fc := GetFileChecker(nil)
		fc.SetExtensions(args)

		if !assert.IsEqual(fc.authorisedExtensions, wantExt) {
			t.Errorf("Set(%+v). Extensions :: Got: %+v. Expected: %+v", args, fc.authorisedExtensions, wantExt)
		}
		if !assert.IsEqual(fc.authorisedTypes, wantTyp) {
			t.Errorf("Set(%+v). Types :: Got: %+v. Expected: %+v", args, fc.authorisedTypes, wantTyp)
		}
	})

	// TEST 3: set the ExtArchiveZIP and ExtAudioMP3 extensions
	// in addition to the default extensions and types set, we expect those
	// two extension to be set as well as the TypeARCHIVE and TypeAUDIO types
	t.Run("ZIP|MP3", func(t *testing.T) {
		var (
			args    = []string{ExtArchiveZIP, ExtAudioMP3}
			wantExt = getExt([]string{ExtArchiveZIP, ExtAudioMP3})
			wantTyp = getTyp([]string{TypeARCHIVE, TypeAUDIO})
		)

		fc := GetFileChecker(nil)
		fc.SetExtensions(args)

		if !assert.IsEqual(fc.authorisedExtensions, wantExt) {
			t.Errorf("Set(%+v). Extensions :: Got: %+v. Expected: %+v", args, fc.authorisedExtensions, wantExt)
		}
		if !assert.IsEqual(fc.authorisedTypes, wantTyp) {
			t.Errorf("Set(%+v). Types :: Got: %+v. Expected: %+v", args, fc.authorisedTypes, wantTyp)
		}
	})
}

func TestFileChecker_UnsetExtensions(t *testing.T) {
	type fields struct {
		file                 *multipart.FileHeader
		authorisedTypes      map[string]bool
		authorisedExtensions map[string]bool
		_dictionary          map[string]string
	}
	type args struct {
		extensions []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FileChecker{
				file:                 tt.fields.file,
				authorisedTypes:      tt.fields.authorisedTypes,
				authorisedExtensions: tt.fields.authorisedExtensions,
				_dictionary:          tt.fields._dictionary,
			}
			fc.UnsetExtensions(tt.args.extensions)
		})
	}
}

func TestFileChecker_IsAuthorised(t *testing.T) {
	body, err := doRequest("assets/fake.png")
	if err != nil {
		t.Error(err)
	}

	t.Log(body)

	// type fields struct {
	// 	file                 *multipart.FileHeader
	// 	authorisedTypes      map[string]bool
	// 	authorisedExtensions map[string]bool
	// 	_dictionary          map[string]string
	// }
	// tests := []struct {
	// 	name   string
	// 	fields fields
	// 	want   bool
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		fc := &FileChecker{
	// 			file:                 tt.fields.file,
	// 			authorisedTypes:      tt.fields.authorisedTypes,
	// 			authorisedExtensions: tt.fields.authorisedExtensions,
	// 			_dictionary:          tt.fields._dictionary,
	// 		}
	// 		if got := fc.IsAuthorised(); got != tt.want {
	// 			t.Errorf("IsAuthorised() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}

//nolint:funlen
func TestFileChecker_isTypeAuthorised(t *testing.T) {
	// default file types authorised: TypeIMAGE & TypeARCHIVE
	type fields struct {
		authorisedTypes map[string]bool
	}
	type args struct {
		header []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "JPG",
			fields: fields{authorisedTypes: getTyp([]string{TypeIMAGE})},
			args:   args{header: []byte{0xFF, 0xD8, 0xFF}},
			want:   true,
		},
		{
			name:   "PNG",
			fields: fields{authorisedTypes: getTyp([]string{TypeIMAGE})},
			args:   args{header: []byte{0x89, 0x50, 0x4E, 0x47}},
			want:   true,
		},
		{
			name:   "PDF",
			fields: fields{authorisedTypes: getTyp([]string{TypeARCHIVE})},
			args:   args{header: []byte{0x25, 0x50, 0x44, 0x46}},
			want:   true,
		},
		{
			name: "XLS-2003",
			// Authorised types: TypeIMAGE, TypeARCHIVE, and TypeDOCUMENTS
			fields: fields{authorisedTypes: getTyp([]string{TypeDOCUMENTS})},
			// Genuine XLS-2003 header
			args: args{header: []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
			want: true,
		},
		{
			name: "XLS-2003-Type-Unauthorised",
			// Authorised types: TypeIMAGE and TypeARCHIVE only
			fields: fields{authorisedTypes: getTyp([]string{})},
			// Genuine XLS-2003 header
			args: args{header: []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}},
			want: false, // because TypeDOCUMENTS not authorised
		},
		{
			name:   "PK-ZIP",
			fields: fields{authorisedTypes: getTyp([]string{TypeARCHIVE})},
			args:   args{header: []byte{0x50, 0x4B, 0x03, 0x04}},
			want:   true,
		},
		{
			name:   "FAKE-PNG",
			fields: fields{authorisedTypes: getTyp([]string{TypeIMAGE})},
			// Fake PNG header
			args: args{header: []byte{0x6E, 0x61, 0x64, 0x69, 0x6D}},
			// though TypeIMAGE among authorised types
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FileChecker{
				authorisedTypes: tt.fields.authorisedTypes,
			}
			if got := fc.isTypeAuthorised(tt.args.header); got != tt.want {
				t.Errorf("isTypeAuthorised() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getExt(setOfExt []string) map[string]bool {
	var extensions = map[string]bool{
		ExtArchivePDF: true,
		ExtImgJPG:     true,
		ExtImgPNG:     true,
	}

	for _, ext := range setOfExt {
		extensions[ext] = true
	}
	return extensions
}

func getTyp(setOfTyp []string) map[string]bool {
	var types = map[string]bool{
		TypeIMAGE:   true,
		TypeARCHIVE: true,
	}

	for _, typ := range setOfTyp {
		types[typ] = true
	}
	return types
}

func mockFileUploadRequest(path string) (*http.Request, error) {
	var (
		err       error
		file      *os.File
		fInfo     os.FileInfo
		ioWriter  io.Writer
		fContents []byte
	)

	if file, err = os.Open(path); err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	if fContents, err = io.ReadAll(file); err != nil {
		return nil, err
	}

	if fInfo, err = file.Stat(); err != nil {
		return nil, err
	}

	var (
		body   = new(bytes.Buffer)
		writer = multipart.NewWriter(body)
	)

	if ioWriter, err = writer.CreateFormFile("uFile", fInfo.Name()); err != nil {
		return nil, err
	}

	_, _ = ioWriter.Write(fContents)
	if err = writer.Close(); err != nil {
		return nil, err
	}

	fmt.Printf("%s", body)
	fmt.Printf("%+v", body)

	return http.NewRequest("POST", "https://google.mu/upload", body)
}

func doRequest(fileName string) ([]byte, error) {
	var (
		err  error
		req  *http.Request
		res  *http.Response
		body []byte

		client = &http.Client{}
		dir, _ = os.Getwd()
		path   = dir + "/" + fileName
	)

	if req, err = mockFileUploadRequest(path); err != nil {
		return nil, err
	}

	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	if _, err = res.Body.Read(body); err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	return body, nil
}
