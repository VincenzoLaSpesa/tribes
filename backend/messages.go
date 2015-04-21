package backend

import (
	"../tools/"
	"bufio"
	"log"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

func NNTP_POST_ReadAndSave(conn net.Conn, groupname string) {

	id_message := tools.RandSeq(32) + "@averno"

	answer_ok := "340 Ok, recommended ID <" + id_message + ">\r\n"
	conn.Write([]byte(answer_ok))
	log.Printf("[FYI] %s", answer_ok)

	var body []string
	var headers []string
	var is_header bool = true

	scanner_h := bufio.NewScanner(conn)
	for {
		ok := scanner_h.Scan()
		if !ok {
			break
		}
		line := scanner_h.Text()

        if strings.HasPrefix(line, "Message-ID:") {continue}
        if (strings.HasPrefix(line, "Newsgroups:")) && (groupname == "garbage") {

                                                    log.Printf("[FYI] ng was %s", groupname)
                                                    grp_hdr:= strings.Split(line,":")
                                                    if grp_hdr != nil {groupname = grp_hdr[1]}
                                                    log.Printf("[FYI] ng now is %s", groupname)
                                                    grp_hdr = strings.Split(groupname,",")
                                                    if grp_hdr != nil {groupname = grp_hdr[0]}
                                                    groupname = strings.Trim(groupname," ")
                                                    log.Printf("[FYI] tried to fix with %s", groupname)

                                                    }

        if strings.HasPrefix(line, "Newsgroups:") {

                                                    line = "Newsgroups: " + groupname
                                                    log.Printf("[FYI] Normalization of NG Header as: %s", line)
                                                    }

        if strings.HasPrefix(line, "Message-ID:") {
                                                    log.Printf("[WARN] not permitted to set MSGID ->%s<-", line)
                                                    continue
                                                    }


		if (line == "") && (is_header == true)  {
			log.Printf("[FYI] header line was ->%s<-", line)
			is_header = false
			continue
		}

		if is_header {
			log.Printf("[FYI] header line is:  ->%s<-", line)
			headers = append(headers, line)
		} else {
			log.Printf("[FYI] body line is:  ->%s<-", line)
			body = append(body, line)

		}
		if line == "." {
			break
		}
	}

	headers = append(headers, "Message-ID: <"+id_message+">")

    num_message, _ := strconv.Atoi(GetNumFilesByGroup(groupname))
	num_message++

    headers = append(headers, "Xref: averno "+groupname+":"+strconv.Itoa(num_message))


	header_file := filepath.Join(messages_folder, "h-"+groupname+"-"+strconv.Itoa(num_message)+"-"+id_message)
	body_file := filepath.Join(messages_folder, "b-"+groupname+"-"+strconv.Itoa(num_message)+"-"+id_message)

	err := tools.WriteMessages(headers, header_file)
	if err != nil {
		log.Printf("[WTF] cannot write headers on %s", header_file)
	} else {
		log.Printf("[FYI] headers saved in %s", header_file)
	}

	err = tools.WriteMessages(body, body_file)
	if err != nil {
		log.Printf("[WTF] cannot write body on %s", body_file)
	} else {
		log.Printf("[FYI] body saved in %s", body_file)
	}

	conn.Write([]byte("240 article posted ok\r\n"))

}

func NNTP_HEAD_ReturnHEADER(conn net.Conn, groupname string, article_id string) {

	article := strings.Trim(article_id, "<")
	article = strings.Trim(article_id, ">")

	if files, err := filepath.Glob(messages_folder + "/h-" + groupname + "-*" + article + "*"); err != nil {
		log.Printf("[SOB] Article %s not found in %s  ", article_id, groupname)
		conn.Write([]byte("430 no such article found\r\n"))
	} else {

		if files == nil {
			log.Printf("[SOB] Article %s not found in %s ", article_id, groupname)
			conn.Write([]byte("430 no such article found\r\n"))
		}

		if files != nil {

            file_fields := strings.Split(files[0],"-")
            head_string := "221 " + file_fields[2] + " <" + file_fields[3] + "> " + "file article retrieved\r\n"
            conn.Write([]byte(head_string))
			Transmit_Article(conn, files[0])
		}

	}
}

func NNTP_BODY_ReturnBODY(conn net.Conn, groupname string, article_id string) {

	article := strings.Trim(article_id, "<")
	article = strings.Trim(article_id, ">")

	if files, err := filepath.Glob(messages_folder + "/b-" + groupname + "-*" + article + "*"); err != nil {
		log.Printf("[SOB] Article %s not found in %s ", article_id, groupname)
		conn.Write([]byte("430 no such article found\r\n"))
	} else {

		if files == nil {
			log.Printf("[SOB] Article %s not found in %s ", article_id, groupname)
			conn.Write([]byte("430 no such article found\r\n"))
		}

		if files != nil {
            file_fields := strings.Split(files[0],"-")
            head_string := "222 " + file_fields[2] + " <" + file_fields[3] + "> " + "file article retrieved\r\n"
            conn.Write([]byte(head_string))
			Transmit_Article(conn, files[0])
		}

	}
}

func NNTP_ARTICLE_ReturnALL(conn net.Conn, groupname string, article_id string) {


    article := strings.Trim(article_id, "<")
	article = strings.Trim(article_id, ">")

    if files, err := filepath.Glob(messages_folder + "/h-" + groupname + "-*" + article + "*"); err != nil {
		log.Printf("[SOB] Article %s not found in %s ", article_id, groupname)
		conn.Write([]byte("430 no such article found\n"))
	} else {

		if files == nil {
			log.Printf("[SOB] Article %s not found in %s ", article_id, groupname)
			conn.Write([]byte("430 no such article found\r\n"))
		}

		if files != nil {
            file_fields := strings.Split(files[0],"-")
            head_string := "220 " + file_fields[2] + " <" + file_fields[3] + "> " + "All of article follows\r\n"
            conn.Write([]byte(head_string))
			Transmit_Article(conn, files[0])
		}

	}




	if files, err := filepath.Glob(messages_folder + "/b-" + groupname + "-*" + article + "*"); err != nil {
		log.Printf("[SOB] Article %s not found in %s ", article_id, groupname)
		conn.Write([]byte("430 no such article found\n"))
	} else {

		if files == nil {
			log.Printf("[SOB] Article %s not found in %s ", article_id, groupname)
			conn.Write([]byte("430 no such article found\r\n"))
		}

		if files != nil {
            head_string := "\r\n\r\n"
            conn.Write([]byte(head_string))
			Transmit_Article(conn, files[0])
		}

	}






}
