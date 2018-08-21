// Package models includes the functions on the model Article.
package models

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

// set flags to output more detailed log
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Article struct {
	Id        int64     `json:"id,omitempty" db:"id" valid:"-"`
	Title     string    `json:"title,omitempty" db:"title" valid:"required,length(10|30)"`
	Text      string    `json:"text,omitempty" db:"text" valid:"required,length(20|4294967295)"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" valid:"-"`
	Comments  []Comment `json:"comments,omitempty" db:"comments" valid:"-"`
}

// DataStruct for the pagination
type ArticlePage struct {
	WhereString string
	WhereParams []interface{}
	Order       map[string]string
	FirstId     int64
	LastId      int64
	PageNum     int
	PerPage     int
	TotalPages  int
	TotalItems  int64
	orderStr    string
}

// Current get the current page of ArticlePage object for pagination.
func (_p *ArticlePage) Current() ([]Article, error) {
	if _, exist := _p.Order["id"]; !exist {
		return nil, errors.New("No id order specified in Order map")
	}
	err := _p.buildPageCount()
	if err != nil {
		return nil, fmt.Errorf("Calculate page count error: %v", err)
	}
	if _p.orderStr == "" {
		_p.buildOrder()
	}
	idStr, idParams := _p.buildIdRestrict("current")
	whereStr := fmt.Sprintf("%s %s %s LIMIT %v", _p.WhereString, idStr, _p.orderStr, _p.PerPage)
	whereParams := []interface{}{}
	whereParams = append(append(whereParams, _p.WhereParams...), idParams...)
	articles, err := FindArticlesWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(articles) != 0 {
		_p.FirstId, _p.LastId = articles[0].Id, articles[len(articles)-1].Id
	}
	return articles, nil
}

// Previous get the previous page of ArticlePage object for pagination.
func (_p *ArticlePage) Previous() ([]Article, error) {
	if _p.PageNum == 0 {
		return nil, errors.New("This's the first page, no previous page yet")
	}
	if _, exist := _p.Order["id"]; !exist {
		return nil, errors.New("No id order specified in Order map")
	}
	err := _p.buildPageCount()
	if err != nil {
		return nil, fmt.Errorf("Calculate page count error: %v", err)
	}
	if _p.orderStr == "" {
		_p.buildOrder()
	}
	idStr, idParams := _p.buildIdRestrict("previous")
	whereStr := fmt.Sprintf("%s %s %s LIMIT %v", _p.WhereString, idStr, _p.orderStr, _p.PerPage)
	whereParams := []interface{}{}
	whereParams = append(append(whereParams, _p.WhereParams...), idParams...)
	articles, err := FindArticlesWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(articles) != 0 {
		_p.FirstId, _p.LastId = articles[0].Id, articles[len(articles)-1].Id
	}
	_p.PageNum -= 1
	return articles, nil
}

// Next get the next page of ArticlePage object for pagination.
func (_p *ArticlePage) Next() ([]Article, error) {
	if _p.PageNum == _p.TotalPages-1 {
		return nil, errors.New("This's the last page, no next page yet")
	}
	if _, exist := _p.Order["id"]; !exist {
		return nil, errors.New("No id order specified in Order map")
	}
	err := _p.buildPageCount()
	if err != nil {
		return nil, fmt.Errorf("Calculate page count error: %v", err)
	}
	if _p.orderStr == "" {
		_p.buildOrder()
	}
	idStr, idParams := _p.buildIdRestrict("next")
	whereStr := fmt.Sprintf("%s %s %s LIMIT %v", _p.WhereString, idStr, _p.orderStr, _p.PerPage)
	whereParams := []interface{}{}
	whereParams = append(append(whereParams, _p.WhereParams...), idParams...)
	articles, err := FindArticlesWhere(whereStr, whereParams...)
	if err != nil {
		return nil, err
	}
	if len(articles) != 0 {
		_p.FirstId, _p.LastId = articles[0].Id, articles[len(articles)-1].Id
	}
	_p.PageNum += 1
	return articles, nil
}

// GetPage is a helper function for the ArticlePage object to return a corresponding page due to
// the parameter passed in, i.e. one of "previous, current or next".
func (_p *ArticlePage) GetPage(direction string) (ps []Article, err error) {
	switch direction {
	case "previous":
		ps, _ = _p.Previous()
	case "next":
		ps, _ = _p.Next()
	case "current":
		ps, _ = _p.Current()
	default:
		return nil, errors.New("Error: wrong dircetion! None of previous, current or next!")
	}
	return
}

// buildOrder is for ArticlePage object to build a SQL ORDER BY clause.
func (_p *ArticlePage) buildOrder() {
	tempList := []string{}
	for k, v := range _p.Order {
		tempList = append(tempList, fmt.Sprintf("%v %v", k, v))
	}
	_p.orderStr = " ORDER BY " + strings.Join(tempList, ", ")
}

// buildIdRestrict is for ArticlePage object to build a SQL clause for ID restriction,
// implementing a simple keyset style pagination.
func (_p *ArticlePage) buildIdRestrict(direction string) (idStr string, idParams []interface{}) {
	switch direction {
	case "previous":
		if strings.ToLower(_p.Order["id"]) == "desc" {
			idStr += "id > ? "
			idParams = append(idParams, _p.FirstId)
		} else {
			idStr += "id < ? "
			idParams = append(idParams, _p.FirstId)
		}
	case "current":
		// trick to make Where function work
		if _p.PageNum == 0 && _p.FirstId == 0 && _p.LastId == 0 {
			idStr += "id > ? "
			idParams = append(idParams, 0)
		} else {
			if strings.ToLower(_p.Order["id"]) == "desc" {
				idStr += "id <= ? AND id >= ? "
				idParams = append(idParams, _p.FirstId, _p.LastId)
			} else {
				idStr += "id >= ? AND id <= ? "
				idParams = append(idParams, _p.FirstId, _p.LastId)
			}
		}
	case "next":
		if strings.ToLower(_p.Order["id"]) == "desc" {
			idStr += "id < ? "
			idParams = append(idParams, _p.LastId)
		} else {
			idStr += "id > ? "
			idParams = append(idParams, _p.LastId)
		}
	}
	if _p.WhereString != "" {
		idStr = " AND " + idStr
	}
	return
}

// buildPageCount calculate the TotalItems/TotalPages for the ArticlePage object.
func (_p *ArticlePage) buildPageCount() error {
	count, err := ArticleCountWhere(_p.WhereString, _p.WhereParams...)
	if err != nil {
		return err
	}
	_p.TotalItems = count
	if _p.PerPage == 0 {
		_p.PerPage = 10
	}
	_p.TotalPages = int(math.Ceil(float64(_p.TotalItems) / float64(_p.PerPage)))
	return nil
}

// FindArticle find a single article by an ID.
func FindArticle(id int64) (*Article, error) {
	if id == 0 {
		return nil, errors.New("Invalid ID: it can't be zero")
	}
	_article := Article{}
	err := DB.Get(&_article, DB.Rebind(`SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles WHERE articles.id = ? LIMIT 1`), id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_article, nil
}

// FirstArticle find the first one article by ID ASC order.
func FirstArticle() (*Article, error) {
	_article := Article{}
	err := DB.Get(&_article, DB.Rebind(`SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles ORDER BY articles.id ASC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_article, nil
}

// FirstArticles find the first N articles by ID ASC order.
func FirstArticles(n uint32) ([]Article, error) {
	_articles := []Article{}
	sql := fmt.Sprintf("SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles ORDER BY articles.id ASC LIMIT %v", n)
	err := DB.Select(&_articles, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _articles, nil
}

// LastArticle find the last one article by ID DESC order.
func LastArticle() (*Article, error) {
	_article := Article{}
	err := DB.Get(&_article, DB.Rebind(`SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles ORDER BY articles.id DESC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_article, nil
}

// LastArticles find the last N articles by ID DESC order.
func LastArticles(n uint32) ([]Article, error) {
	_articles := []Article{}
	sql := fmt.Sprintf("SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles ORDER BY articles.id DESC LIMIT %v", n)
	err := DB.Select(&_articles, DB.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _articles, nil
}

// FindArticles find one or more articles by the given ID(s).
func FindArticles(ids ...int64) ([]Article, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return nil, errors.New(msg)
	}
	_articles := []Article{}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := DB.Rebind(fmt.Sprintf(`SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles WHERE articles.id IN (?%s)`, idsHolder))
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	err := DB.Select(&_articles, sql, idsT...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _articles, nil
}

// FindArticleBy find a single article by a field name and a value.
func FindArticleBy(field string, val interface{}) (*Article, error) {
	_article := Article{}
	sqlFmt := `SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles WHERE %s = ? LIMIT 1`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err := DB.Get(&_article, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_article, nil
}

// FindArticlesBy find all articles by a field name and a value.
func FindArticlesBy(field string, val interface{}) (_articles []Article, err error) {
	sqlFmt := `SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles WHERE %s = ?`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err = DB.Select(&_articles, DB.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _articles, nil
}

// AllArticles get all the Article records.
func AllArticles() (articles []Article, err error) {
	err = DB.Select(&articles, "SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return articles, nil
}

// ArticleCount get the count of all the Article records.
func ArticleCount() (c int64, err error) {
	err = DB.Get(&c, "SELECT count(*) FROM articles")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// ArticleCountWhere get the count of all the Article records with a where clause.
func ArticleCountWhere(where string, args ...interface{}) (c int64, err error) {
	sql := "SELECT count(*) FROM articles"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	err = stmt.Get(&c, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// ArticleIncludesWhere get the Article associated models records, currently it's not same as the corresponding "includes" function but "preload" instead in Ruby on Rails. It means that the "sql" should be restricted on Article model.
func ArticleIncludesWhere(assocs []string, sql string, args ...interface{}) (_articles []Article, err error) {
	_articles, err = FindArticlesWhere(sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(assocs) == 0 {
		log.Println("No associated fields ard specified")
		return _articles, err
	}
	if len(_articles) <= 0 {
		return nil, errors.New("No results available")
	}
	ids := make([]interface{}, len(_articles))
	for _, v := range _articles {
		ids = append(ids, interface{}(v.Id))
	}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	for _, assoc := range assocs {
		switch assoc {
		case "comments":
			where := fmt.Sprintf("article_id IN (?%s)", idsHolder)
			_comments, err := FindCommentsWhere(where, ids...)
			if err != nil {
				log.Printf("Error when query associated objects: %v\n", assoc)
				continue
			}
			for _, vv := range _comments {
				for i, vvv := range _articles {
					if vv.ArticleId == vvv.Id {
						vvv.Comments = append(vvv.Comments, vv)
					}
					_articles[i].Comments = vvv.Comments
				}
			}
		}
	}
	return _articles, nil
}

// ArticleIds get all the IDs of Article records.
func ArticleIds() (ids []int64, err error) {
	err = DB.Select(&ids, "SELECT id FROM articles")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}

// ArticleIdsWhere get all the IDs of Article records by where restriction.
func ArticleIdsWhere(where string, args ...interface{}) ([]int64, error) {
	ids, err := ArticleIntCol("id", where, args...)
	return ids, err
}

// ArticleIntCol get some int64 typed column of Article by where restriction.
func ArticleIntCol(col, where string, args ...interface{}) (intColRecs []int64, err error) {
	sql := "SELECT " + col + " FROM articles"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&intColRecs, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return intColRecs, nil
}

// ArticleStrCol get some string typed column of Article by where restriction.
func ArticleStrCol(col, where string, args ...interface{}) (strColRecs []string, err error) {
	sql := "SELECT " + col + " FROM articles"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&strColRecs, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strColRecs, nil
}

// FindArticlesWhere query use a partial SQL clause that usually following after WHERE
// with placeholders, eg: FindUsersWhere("first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18.
func FindArticlesWhere(where string, args ...interface{}) (articles []Article, err error) {
	sql := "SELECT COALESCE(articles.text, '') AS text, articles.id, articles.title, articles.created_at, articles.updated_at FROM articles"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&articles, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return articles, nil
}

// FindArticleBySql query use a complete SQL clause
// with placeholders, eg: FindUserBySql("SELECT * FROM users WHERE first_name = ? AND age > ? ORDER BY DESC LIMIT 1", "John", 18)
// will return only One record in the table "users" whose first_name is "John" and age elder than 18.
func FindArticleBySql(sql string, args ...interface{}) (*Article, error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_article := &Article{}
	err = stmt.Get(_article, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return _article, nil
}

// FindArticlesBySql query use a complete SQL clause
// with placeholders, eg: FindUsersBySql("SELECT * FROM users WHERE first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18.
func FindArticlesBySql(sql string, args ...interface{}) (articles []Article, err error) {
	stmt, err := DB.Preparex(DB.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&articles, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return articles, nil
}

// CreateArticle use a named params to create a single Article record.
// A named params is key-value map like map[string]interface{}{"first_name": "John", "age": 23} .
func CreateArticle(am map[string]interface{}) (int64, error) {
	if len(am) == 0 {
		return 0, fmt.Errorf("Zero key in the attributes map!")
	}
	t := time.Now()
	for _, v := range []string{"created_at", "updated_at"} {
		if am[v] == nil {
			am[v] = t
		}
	}
	keys := allKeys(am)
	sqlFmt := `INSERT INTO articles (%s) VALUES (%s)`
	sql := fmt.Sprintf(sqlFmt, strings.Join(keys, ","), ":"+strings.Join(keys, ",:"))
	result, err := DB.NamedExec(sql, am)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return lastId, nil
}

// Create is a method for Article to create a record.
func (_article *Article) Create() (int64, error) {
	ok, err := govalidator.ValidateStruct(_article)
	if !ok {
		errMsg := "Validate Article struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Article struct error: " + err.Error()
		}
		log.Println(errMsg)
		return 0, errors.New(errMsg)
	}
	t := time.Now()
	_article.CreatedAt = t
	_article.UpdatedAt = t
	sql := `INSERT INTO articles (title,text,created_at,updated_at) VALUES (:title,:text,:created_at,:updated_at)`
	result, err := DB.NamedExec(sql, _article)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return lastId, nil
}

// CommentsCreate is used for Article to create the associated objects Comments
func (_article *Article) CommentsCreate(am map[string]interface{}) error {
	am["article_id"] = _article.Id
	_, err := CreateComment(am)
	return err
}

// GetComments is used for Article to get associated objects Comments
// Say you have a Article object named article, when you call article.GetComments(),
// the object will get the associated Comments attributes evaluated in the struct.
func (_article *Article) GetComments() error {
	_comments, err := ArticleGetComments(_article.Id)
	if err == nil {
		_article.Comments = _comments
	}
	return err
}

// ArticleGetComments a helper fuction used to get associated objects for ArticleIncludesWhere().
func ArticleGetComments(id int64) ([]Comment, error) {
	_comments, err := FindCommentsBy("article_id", id)
	return _comments, err
}

// Destroy is method used for a Article object to be destroyed.
func (_article *Article) Destroy() error {
	if _article.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := DestroyArticle(_article.Id)
	return err
}

// DestroyArticle will destroy a Article record specified by the id parameter.
func DestroyArticle(id int64) error {
	// Destroy association objects at first
	// Not care if exec properly temporarily
	destroyArticleAssociations(id)
	stmt, err := DB.Preparex(DB.Rebind(`DELETE FROM articles WHERE id = ?`))
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

// DestroyArticles will destroy Article records those specified by the ids parameters.
func DestroyArticles(ids ...int64) (int64, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return 0, errors.New(msg)
	}
	// Destroy association objects at first
	// Not care if exec properly temporarily
	destroyArticleAssociations(ids...)
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := fmt.Sprintf(`DELETE FROM articles WHERE id IN (?%s)`, idsHolder)
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	result, err := stmt.Exec(idsT...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// DestroyArticlesWhere delete records by a where clause restriction.
// e.g. DestroyArticlesWhere("name = ?", "John")
// And this func will not call the association dependent action
func DestroyArticlesWhere(where string, args ...interface{}) (int64, error) {
	sql := `DELETE FROM articles WHERE `
	if len(where) > 0 {
		sql = sql + where
	} else {
		return 0, errors.New("No WHERE conditions provided")
	}
	ids, x_err := ArticleIdsWhere(where, args...)
	if x_err != nil {
		log.Printf("Delete associated objects error: %v\n", x_err)
	} else {
		destroyArticleAssociations(ids...)
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// destroyArticleAssociations is a private function used to destroy a Article record's associated objects.
// The func not return err temporarily.
func destroyArticleAssociations(ids ...int64) {
	idsHolder := ""
	if len(ids) > 1 {
		idsHolder = strings.Repeat(",?", len(ids)-1)
	}
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	var err error
	// make sure no declared-and-not-used exception
	_, _, _ = idsHolder, idsT, err
	where := fmt.Sprintf("article_id IN (?%s)", idsHolder)
	_, err = DestroyCommentsWhere(where, idsT...)
	if err != nil {
		log.Printf("Destroy associated object %s error: %v\n", "Comments", err)
	}
}

// Save method is used for a Article object to update an existed record mainly.
// If no id provided a new record will be created. FIXME: A UPSERT action will be implemented further.
func (_article *Article) Save() error {
	ok, err := govalidator.ValidateStruct(_article)
	if !ok {
		errMsg := "Validate Article struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Article struct error: " + err.Error()
		}
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if _article.Id == 0 {
		_, err = _article.Create()
		return err
	}
	_article.UpdatedAt = time.Now()
	sqlFmt := `UPDATE articles SET %s WHERE id = %v`
	sqlStr := fmt.Sprintf(sqlFmt, "title = :title, text = :text, updated_at = :updated_at", _article.Id)
	_, err = DB.NamedExec(sqlStr, _article)
	return err
}

// UpdateArticle is used to update a record with a id and map[string]interface{} typed key-value parameters.
func UpdateArticle(id int64, am map[string]interface{}) error {
	if len(am) == 0 {
		return errors.New("Zero key in the attributes map!")
	}
	am["updated_at"] = time.Now()
	keys := allKeys(am)
	sqlFmt := `UPDATE articles SET %s WHERE id = %v`
	setKeysArr := []string{}
	for _, v := range keys {
		s := fmt.Sprintf(" %s = :%s", v, v)
		setKeysArr = append(setKeysArr, s)
	}
	sqlStr := fmt.Sprintf(sqlFmt, strings.Join(setKeysArr, ", "), id)
	_, err := DB.NamedExec(sqlStr, am)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Update is a method used to update a Article record with the map[string]interface{} typed key-value parameters.
func (_article *Article) Update(am map[string]interface{}) error {
	if _article.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateArticle(_article.Id, am)
	return err
}

// UpdateAttributes method is supposed to be used to update Article records as corresponding update_attributes in Ruby on Rails.
func (_article *Article) UpdateAttributes(am map[string]interface{}) error {
	if _article.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateArticle(_article.Id, am)
	return err
}

// UpdateColumns method is supposed to be used to update Article records as corresponding update_columns in Ruby on Rails.
func (_article *Article) UpdateColumns(am map[string]interface{}) error {
	if _article.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateArticle(_article.Id, am)
	return err
}

// UpdateArticlesBySql is used to update Article records by a SQL clause
// using the '?' binding syntax.
func UpdateArticlesBySql(sql string, args ...interface{}) (int64, error) {
	if sql == "" {
		return 0, errors.New("A blank SQL clause")
	}
	stmt, err := DB.Preparex(DB.Rebind(sql))
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
