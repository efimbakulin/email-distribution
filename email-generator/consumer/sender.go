package main

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

var (
	connCache = make(map[string]*smtp.Client)
)

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func dialExchange(srv string) (*smtp.Client, error) {
	if connCache[srv] != nil {
		return connCache[srv], nil
	}
	conn, err := smtp.Dial(srv)
	if err != nil {
		return nil, err
	}
	connCache[srv] = conn
	return conn, nil
}

func sendMail(subj string, body string, from string, fromName string, to string) error {
	// Set up authentication information.

	smtpServer := "mail.ocplay.net:25"

	conn, err := dialExchange(smtpServer)

	if err != nil {
		return err
	}

	if err = conn.Mail(from); err != nil {
		return err
	}

	if err = conn.Rcpt(to); err != nil {
		return err
	}

	writer, err := conn.Data()
	if err != nil {
		return err
	}
	fromAddr := mail.Address{fromName, from}

	header := make(map[string]string)
	header["From"] = fromAddr.String()
	header["To"] = to
	header["Subject"] = encodeRFC2047(subj)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	for k, v := range header {
		writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
	}
	writer.Write([]byte("\r\n"))
	writer.Write([]byte(base64.StdEncoding.EncodeToString([]byte(body))))
	writer.Close()

	return nil
}
