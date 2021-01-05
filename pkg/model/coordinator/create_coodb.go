package coordinator

import (
	// ビルド時のみ使用する
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB Path(相対パスでも大丈夫かと思うが、筆者の場合、絶対パスでないと実行できなかった)
const dbPath = "/home/mash/hornet/coodb.sql"

// コネクションプールを作成
var DbConnection *sql.DB

// データ格納用
type Coomile struct {
	index int
	tag   string
}

func create_coodb(index int, txTag string) {
	// Open(driver,  sql 名(任意の名前))
	DbConnection, _ := sql.Open("sqlite3", dbPath)

	// Connection をクローズする。(defer で閉じるのが Golang の作法)
	defer DbConnection.Close()

	// blog テーブルの作成
	cmd := `CREATE TABLE IF NOT EXISTS coomile(
             index INT,        
             tag STRING)`

	// cmd を実行
	// _ -> 受け取った結果に対して何もしないので、_ にする
	_, err := DbConnection.Exec(cmd)

	// エラーハンドリング(Go だと大体このパターン)
	if err != nil {
		// Fatalln は便利
		// エラーが発生した場合、以降の処理を実行しない
		log.Fatalln(err)
	}

	cmd = "INSERT INTO coomile (index, tag) VALUES (?, ?)"
	_, err = DbConnection.Exec(cmd, index, txTag)

	if err != nil {
		// golang には、try-catch がない。nil か否かで判定
		log.Fatalln(err)
	}

	/*ここから挿入したデータの一覧を出力する処理*/
	// マルチプルセレクト(今度は、_ ではなく、rows)
	cmd = "SELECT * FROM coomile where index = ?"
	row := DbConnection.QueryRow(cmd, index)

	// データ保存領域を確保
	var b Coomile
	// Scan にて、struct のアドレスにデータを入れる
	err = row.Scan(&b.index, &b.tag)
	// エラーハンドリング(共通関数にした方がいいのかな)
	if err != nil {
		// シングルセレクトの場合は、エラーハンドリングが異なる
		if err == sql.ErrNoRows {
			log.Println("There is no row!!!")
		} else {
			log.Println(err)
		}
	}
	fmt.Println(b.index, b.tag)

	log.Println("create_coodb")
}