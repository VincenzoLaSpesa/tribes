package nntp

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
	"tribes/3be"
	"tribes/backend"
)

var capab_out string = "101 Capability list:\nVERSION 2\nREADER\nPOST\nSTAT\nXOVER\nOVER\nLIST ACTIVE NEWSGROUPS OVERVIEW.FMT\n"

func init() {

	go NNTP_Frontend()

}

func NNTP_Frontend() {

	// setting up the tcp connection

	ln, err := net.Listen("tcp", "127.0.0.1:11119")
	if err == nil {
		log.Printf("[NNTP] TCP listening at %s ", "127.0.0.1:11119")

	} else {
		log.Printf("[WTF] TCP CANNOT listen at %s. SYSADMIIIIN!!", "127.0.0.1:11119")
		return
	}

	defer ln.Close()

	for {

		// start listening at it

		server, err := ln.Accept()
		tcp_client := server.RemoteAddr()

		if err == nil {

			log.Printf("[INFO] NNTP accepted connection from %s ", tcp_client)

		} else {
			log.Printf("[WTF] NNTP something went wrong at %s. SYSADMIIIIN!!", "127.0.0.1:11119")
		}

		// start the NNTP interpreter in background.
		go NNTP_Interpret(server)

	}

}

func NNTP_Interpret(conn net.Conn) {

	var current_group string = "garbage"
	var current_messg string = "null"

	remote_client := conn.RemoteAddr()
	greetings := "200 averno.node AVERNO Version 01 beta, S0, posting OK"
	conn.Write([]byte(greetings + "\r\n"))
	for {
		linea := make([]byte, 1024)

		linea, _, _ = bufio.NewReader(conn).ReadLine()

		message := string(linea)

		// decides WTF to do with the string

		if matches, _ := regexp.MatchString("(?i)^QUIT.*", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			conn.Write([]byte("205 closing connection - goodbye!"))
			conn.Close()
			break
		}

		if matches, _ := regexp.MatchString("(?i)^GROUP[ ]+([0-9A-Za-z]+\\.)+[0-9A-Za-z]+$", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			sinta := strings.Split(message, " ")
			current_group = sinta[1]
			tribe.BCastGroup(current_group)
			conn.Write([]byte(backend.ResponseToNNTPGROUP(current_group)))
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^GROUP[ ]*$", message); matches == true {
			log.Printf("[WTF] NNTP %s from %s with no groupname ", message, remote_client)
			antwort := "501 Syntax is:  GROUP newsgroup\r\n"
			conn.Write([]byte(antwort))
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^STAT[ ]+.+", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			sinta := strings.Split(message, " ")
			current_messg = sinta[1]
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^LIST[ ]*$", message); matches == true {
			conn.Write([]byte("215 list of newsgroups follows\r\n"))
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			backend.Trasmit_Active_NG(conn)
			conn.Write([]byte(".\r\n"))
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^LIST[ ]+OVERVIEW.FMT[ ]*$", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			backend.NNTP_LIST_OVERVIEW_FMT(conn)
			continue
		}

		// split. To consolidate later

		if matches, _ := regexp.MatchString("(?i)^HEAD[ ]*$", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			backend.NNTP_HEAD_ReturnHEADER(conn, current_group, current_messg)
			conn.Write([]byte(".\r\n"))
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^HEAD[ ](([0-9]+)|(<.*>))", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			sinta := strings.Split(message, " ")
			current_messg = sinta[1]
			backend.NNTP_HEAD_ReturnHEADER(conn, current_group, current_messg)
			conn.Write([]byte(".\r\n"))
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^BODY[ ]*$", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			backend.NNTP_BODY_ReturnBODY(conn, current_group, current_messg)
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^BODY[ ](([0-9]+)|(<.*>))", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			sinta := strings.Split(message, " ")
			current_messg = sinta[1]
			backend.NNTP_BODY_ReturnBODY(conn, current_group, current_messg)
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^ARTICLE[ ]*$", message); matches == true {
			backend.NNTP_ARTICLE_ReturnALL(conn, current_group, current_messg)
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^ARTICLE[ ](([0-9]+)|(<.*>))", message); matches == true {
			sinta := strings.Split(message, " ")
			current_messg = sinta[1]
			backend.NNTP_ARTICLE_ReturnALL(conn, current_group, current_messg)
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^POST[ ]*$", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			backend.NNTP_POST_ReadAndSave(conn, current_group)
			continue
		}
		if matches, _ := regexp.MatchString("(?i)^STAT[ ](([0-9]+)|(<.*>))", message); matches == true {
			sinta := strings.Split(message, " ")
			current_messg = sinta[1]
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			continue
		}
		if matches, _ := regexp.MatchString("(?i)^CAPABILITIES[ ]*$", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			conn.Write([]byte(capab_out))
			continue

		}
		if matches, _ := regexp.MatchString("(?i)^MODE[ ]*READER[ ]*$", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			conn.Write([]byte("200 Hello, you can post\r\n"))
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^NEWGROUPS.*", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			conn.Write([]byte("231 New newsgroups since whatever follow" + "\r\n"))
			backend.Trasmit_New_NG(conn)
			conn.Write([]byte("\r\n." + "\r\n"))
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^X?OVER[ ]+.+", message); matches == true {
			log.Printf("[INFO] NNTP %s from %s ", message, remote_client)
			backend.NNTP_XOVER_Frontend(current_group, conn, message)
			continue
		}

		if matches, _ := regexp.MatchString("(?i)^X?OVER[ ]*", message); matches == true {
			log.Printf("[WTF] NNTP %s from %s (with NO group!)", message, remote_client)
			conn.Write([]byte("412 Not in a newsgroup\n\r"))
			continue
		}

		if message == "" {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("[INFO] NNTP BULLSHIT >%s< from %s ", message, remote_client)
		time.Sleep(1 * time.Second)
		conn.Write([]byte("500 Command not understood\r\n"))

	}
	conn.Close()
}

func NNTP_Engine_Start() {

	log.Printf("[NNTP] NNTP Engine Starting")

}
