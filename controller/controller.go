package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/SeijiOmi/posts-service/entity"
	"github.com/SeijiOmi/posts-service/service"
)

// Index action: GET /posts
func Index(c *gin.Context) {
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		return
	}
	var b service.Behavior
	p, err := b.GetAllAttachJoinData(offset)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// Create action: POST /posts
func Create(c *gin.Context) {
	var inputPost entity.Post
	if err := bindJSON(c, &inputPost); err != nil {
		return
	}
	var inputJoinPost entity.JoinPost
	if err := bindJSON(c, &inputJoinPost); err != nil {
		return
	}
	inputJoinPost.Post = inputPost

	type tokenStru struct {
		Token string `json:"token"`
	}
	var token tokenStru
	if err := bindJSON(c, &token); err != nil {
		return
	}

	var b service.Behavior
	createdPost, err := b.CreateModel(inputJoinPost, token.Token)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, createdPost)
	}
}

// Show action: GET /posts/:id
func Show(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior
	p, err := b.GetByID(id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// Update action: PUT /posts/:id
func Update(c *gin.Context) {
	id := c.Params.ByName("id")
	var inputPost entity.Post
	if err := bindJSON(c, &inputPost); err != nil {
		return
	}

	var b service.Behavior
	p, err := b.UpdateByID(id, inputPost)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, p)
	}
}

// Delete action: DELETE /posts/:id
func Delete(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior

	if err := b.DeleteByID(id); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, gin.H{"id #" + id: "deleted"})
	}
}

// UserShow action: get /user/:id
func UserShow(c *gin.Context) {
	id := c.Params.ByName("id")
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
		return
	}

	var b service.Behavior
	p, err := b.GetByUserIDAttachJoinData(id, offset)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// HelperShow action: get /helpser/:id
func HelperShow(c *gin.Context) {
	id := c.Params.ByName("id")
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
		return
	}

	var b service.Behavior
	p, err := b.GetByHelperUserIDAttachJoinData(id, offset)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// SetHelpUser action: Post /helper
func SetHelpUser(c *gin.Context) {
	id, token, err := bindGetIDAndToken(c)
	fmt.Println(id)
	if err != nil {
		return
	}

	var b service.Behavior
	p, err := b.SetHelpUserID(id, token)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, p)
	}
}

// TakeHelpUser action: delete /helper
func TakeHelpUser(c *gin.Context) {
	id := c.Params.ByName("id")
	_, token, err := bindGetIDAndToken(c)
	if err != nil {
		return
	}

	var b service.Behavior
	p, err := b.TakeHelpUserID(id, token)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, p)
	}
}

// DonePayment action: POST /done
func DonePayment(c *gin.Context) {
	id, token, err := bindGetIDAndToken(c)
	if err != nil {
		return
	}

	var b service.Behavior
	p, err := b.DonePayment(id, token)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, p)
	}
}

// DoneAcceptance action: PUT /done
func DoneAcceptance(c *gin.Context) {
	id := c.Params.ByName("id")
	_, token, err := bindGetIDAndToken(c)
	if err != nil {
		return
	}

	var b service.Behavior
	p, err := b.DoneAcceptance(id, token)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, p)
	}
}

// AmountPayment action: GET /amount-payment
func AmountPayment(c *gin.Context) {
	id := c.Params.ByName("id")

	var b service.Behavior
	p, err := b.GetAmountPaymentByUserID(id)

	response := struct {
		AmountPayment int
	}{
		p,
	}

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

// TagShow action: GET /tag/id
func TagShow(c *gin.Context) {
	id := c.Params.ByName("id")
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
		return
	}

	var b service.Behavior
	p, err := b.GetByTagIDAttachJoinData(id, offset)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// TagLike action: GET /tag/like
func TagLike(c *gin.Context) {
	id := c.Params.ByName("id")

	var b service.Behavior
	p, err := b.FindTagLikeBody(id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, p)
	}
}

func bindGetIDAndToken(c *gin.Context) (string, string, error) {
	type requestStru struct {
		ID    float64 `json:"id"`
		Token string  `json:"token"`
	}
	var request requestStru
	if err := bindJSON(c, &request); err != nil {
		return "", "", err
	}

	return strconv.Itoa(int(request.ID)), request.Token, nil
}

func bindJSON(c *gin.Context, data interface{}) error {
	buf := make([]byte, 2048)
	n, _ := c.Request.Body.Read(buf)
	b := string(buf[0:n])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(b)))
	if err := c.BindJSON(data); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println("bind JSON err")
		fmt.Println(err)
		return err
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(b)))
	return nil
}
