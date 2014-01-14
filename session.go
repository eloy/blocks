package blocks

import (
	"net/http"
	"encoding/json"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
)

type SessionManager interface {
	Set(string, string)
	Get(string) string
}

type sessionManager struct {
	request *Request
	session map[string]string
}

func newSessionManager(r *Request) *sessionManager {
	s := new(sessionManager)
	s.request = r
	s.session = make(map[string]string, 0)
	return s
}

func (this *sessionManager) read() {
	content, found := this.readCookieContent()
	if !found {
		return
	}

	err := json.Unmarshal(content, &this.session)
	if err != nil {
		this.session = make(map[string]string, 0)
	}
}


func (this *sessionManager) save() {
	content, _ := json.Marshal(this.session)
	this.writeCookieContent(content )
}

func (this *sessionManager) Set(key string, value string) {
	this.session[key] = value
}

func (this *sessionManager) Get(key string) string {
	return this.session[key]
}


func (this *sessionManager) readCookieContent() ([]byte, bool) {
	c, err := this.request.serverRequest.Cookie(sessionCookieName());
	if err == nil {
		return decrypt(c.Value)
	} else if err == http.ErrNoCookie {
		return nil, false
	} else {
		panic(err)
	}
}

func (this *sessionManager) writeCookieContent(content []byte) {

	cookie := &http.Cookie{
		Name: sessionCookieName(),
		Value: encrypt(content),
	}

	http.SetCookie(this.request.writer, cookie)
}


func sessionCookieName() string {
	return "blocks_session"
}


var cookieCipher = createCookieCipher()
func createCookieCipher() cipher.Block {
	key := []byte {166, 67, 129, 54, 165, 157, 195, 196, 228, 63, 218, 48, 188, 201, 230, 216}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}


// encrypt string to base64 crypto using AES
func encrypt(plaintext []byte) string {
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(cookieCipher, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decrypt(cryptoText string) ([]byte, bool) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		log.Println("ERROR: ciphertext too short")
		return nil, false
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(cookieCipher, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, true
}
