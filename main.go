package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"time"
)
/**
Specifications of the rectangle
 */
type InputRect struct {
	X *int `json:"x" form:"x" binding:"required"`
	Y *int `json:"y" form:"y" binding:"required"`
	Width *int `json:"width" binding:"required"`
	Height *int `json:"height" binding:"required"`
}
/**
input Json format
 */
type InputData struct {
	Main    *InputRect `json:"main" form:"main" binding:"required"`
	Input   []*InputRect `json:"input" binding:"required"`
}
/**
Main Json model, this one contains gorm.Model which adds ID, Created_at, etc. however we use time for
the sake of having control over our values and in case plugin changes its behavior or we decide to use another
plugin.
 */
type DbModel struct {
	gorm.Model
	Time time.Time `gorm:"column:time" time_format:"2006-01-02 15:04:05" `
	X int `gorm:"column:x"`
	Y int `gorm:"column:y"`
	Width int
	Height int
}
/**
This struct is to make our Json transform closest to the format we want. Doesn't contain plugin's default model.
time_format doesn't quite apply to gorm and it needs to be implemented in frontend.
 */
type DbModelFind struct {
	Time time.Time `gorm:"column:time" time_format:"2006-01-02 15:04:05"`
	X int `gorm:"column:x"`
	Y int `gorm:"column:y"`
	Width int
	Height int
}

func main() {
	/**
	DB init.
	 */
	db, err := gorm.Open("sqlite3", "inputs.db")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&DbModel{})
	/**
	Gin init.
	 */
	r := gin.Default()
	/**
	Implements POST
	*/
	r.POST("/", func(c *gin.Context)  {
		var input InputData
		//Check if Payload is JSON and conforms to our format
		if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//i :=0
		for _, s := range input.Input {
			if intersects(input.Main, s) {
				//input.Input[i] = s
				/**
				In case this rectangle has some area in common with Main, lets save it.
				 */
				db.Create(&DbModel{
					//Time: time.Now().Format("2006-01-02 15:04:05"),
					Time: time.Now(),
					X: *s.X,
					Y: *s.Y,
					Width: *s.Width,
					Height: *s.Height})
			}
		}
		//The commented code bellow is left intentionally for reference purposes.
		/*input.Input = input.Input[:i]
		b, err := json.Marshal(input)
		if err != nil {
			fmt.Println(err)
			return
		}
		db.Create(&DbModel{Time: time.Now(), Data: string(b)})*/
		c.JSON(http.StatusOK, gin.H{"message": "Your input Was Parsed Successfully."})
	})

	/**
	Implements GET
	 */
	r.GET("/", func(c *gin.Context)  {
		Everything := []DbModelFind {}

		db.Table("db_models").Find(&Everything)
		//Show 'em wazzup!
		c.JSON(http.StatusOK, Everything)
	})
	r.Run(":8080")
}

func intersects(rect1 *InputRect, rect2 *InputRect) bool {
	//if:
	DoesNotOverlap :=
		//start X of rect1 is bigger than max X of rect2, rect2 is in right of rect1
		*rect1.X > *rect2.X+*rect2.Width ||
		//start X of rect2 is bigger than max X of rect1, rect1 is in right of rect2
		*rect2.X > *rect1.X+*rect1.Width ||
		//start Y of rect1 is bigger than max Y of rect2, rect1 is on top of rect1
		*rect1.Y > *rect2.Y+*rect2.Height ||
		//start Y of rect2 is bigger than max Y of rect1, rect2 is on top of rect2
		*rect2.Y > *rect1.Y+*rect1.Height
	/**
	DoesNotOverlap: there are no case in which rect1 and rect2 intersect.
	 */

	return !DoesNotOverlap // reverse it and shows whether or not it intersects.


}

//TODO: Write Testing.