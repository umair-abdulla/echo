package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {

	dbase, err := sql.Open("postgres", "user=postgres password=abdullaumair dbname=db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err = dbase.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("DB Connected...")
	}
	e := echo.New()

	type Student struct {
		StudentId string `json:"student_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       string `json:"age"`
	}
	type students struct {
		Students []Student `json:"students"`
	}

	type course struct {
		CourseId   string `json:"course_id"`
		CourseName string `json:"course_name"`
		CourseFee  string `json:"course_fee"`
	}
	type courses struct {
		Courses []course `json:"courses"`
	}

	type registration struct {
		StudentId string `json:"student_id"`
		CoursesId string `json:"courses_id"`
		RegDate   string `json:"reg_date"`
	}
	type registrations struct {
		registrations []registration `json:"registrations"`
	}

	e.GET("/student", func(c echo.Context) error {
		sqlStatement := "SELECT student_id, first_name, last_name, age FROM public.student order by student_id"
		rows, err := dbase.Query(sqlStatement)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		}
		defer rows.Close()
		result := students{Students: make([]Student, 0, 10)}
		i := 0
		for rows.Next() {
			student := Student{}

			err2 := rows.Scan(&student.StudentId, &student.FirstName, &student.LastName, &student.Age)

			// Exit if we get an error
			if err2 != nil {
				return err2
			}
			result.Students = append(result.Students, student)
			i++
		}
		fmt.Print(result.Students[0])
		return c.JSON(http.StatusCreated, result.Students)

		//return c.String(http.StatusOK, "ok")
	})
	e.GET("/student/:studentid", func(c echo.Context) error {
		studentId := c.Param("studentid")
		sqlStatement := "SELECT student_id, first_name, last_name, age FROM public.student where student_id=$1"
		row := dbase.QueryRow(sqlStatement, studentId)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		}

		student := Student{}
		err = row.Scan(&student.StudentId, &student.FirstName, &student.LastName, &student.Age)
		if err != nil {
			return c.JSON(http.StatusNotFound, "")
		}

		return c.JSON(http.StatusCreated, student)

		//return c.String(http.StatusOK, "ok")
	})
	e.GET("/course", func(c echo.Context) error {
		sqlStatement := "SELECT course_id, course_name, course_fee FROM public.course order by course_id"
		rows, err := dbase.Query(sqlStatement)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		}
		defer rows.Close()
		result := courses{Courses: make([]course, 0, 10)}
		i := 0
		for rows.Next() {
			course := course{}

			err2 := rows.Scan(&course.CourseId, &course.CourseName, &course.CourseFee)

			// Exit if we get an error
			if err2 != nil {
				return err2
			}
			result.Courses = append(result.Courses, course)
			i++
		}
		fmt.Print(result.Courses[0])
		return c.JSON(http.StatusCreated, result.Courses)

		//return c.String(http.StatusOK, "ok")
	})
	e.GET("/course/:courseid", func(c echo.Context) error {
		courseId := c.Param("courseid")
		sqlStatement := "SELECT course_id, course_name, course_fee FROM public.course where course_id=$1"
		row := dbase.QueryRow(sqlStatement, courseId)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		}

		course := course{}
		err = row.Scan(&course.CourseId, &course.CourseName, &course.CourseFee)
		if err != nil {
			return c.JSON(http.StatusNotFound, "")
		}

		return c.JSON(http.StatusCreated, course)
	})

	e.POST("/student", func(c echo.Context) error {
		u := new(Student)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlStatement := "INSERT INTO public.\"student\" (student_id, first_name, last_name, age)VALUES ($1, $2, $3, $4)"
		res, err := dbase.Query(sqlStatement, u.StudentId, u.FirstName, u.LastName, u.Age)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, "ok")
	})

	e.PUT("/student/:studentid", func(c echo.Context) error {
		studentid := c.Param("studentid")
		u := new(Student)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlStatement := "UPDATE student SET first_name=$1,last_name=$2, age=$3 WHERE student_id=$4"
		res, err := dbase.Query(sqlStatement, u.FirstName, u.LastName, u.Age, studentid)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, u.StudentId)
	})

	e.GET("/reg", func(c echo.Context) error {
		sqlStatement := "SELECT stud_id, courses_id, reg_date FROM public.reg order by courses_id"
		rows, err := dbase.Query(sqlStatement)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		}
		defer rows.Close()
		result := registrations{registrations: make([]registration, 0, 10)}
		i := 0
		for rows.Next() {
			registration := registration{}

			err2 := rows.Scan(&registration.CoursesId, &registration.StudentId, &registration.RegDate)

			// Exit if we get an error
			if err2 != nil {
				return err2
			}
			result.registrations = append(result.registrations, registration)
			i++
		}
		fmt.Print(result.registrations[0])
		return c.JSON(http.StatusCreated, result.registrations)
	})
	e.GET("/course/:courseid/reg", func(c echo.Context) error {
		courseId := c.Param("courseid")
		sqlStatement := "SELECT courses_id, stud_id, reg_date FROM public.reg where courses_id=$1 order by stud_id"
		row := dbase.QueryRow(sqlStatement, courseId)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		}

		registration := registration{}
		err = row.Scan(&registration.CoursesId, &registration.StudentId, &registration.RegDate)
		if err != nil {
			return c.JSON(http.StatusNotFound, "")
		}

		return c.JSON(http.StatusCreated, registration)
	})

	e.GET("/student/:studentid/reg", func(c echo.Context) error {
		studentId := c.Param("studentid")
		sqlStatement := "SELECT courses_id, stud_id, reg_date FROM public.reg where stud_id=$1 order by stud_id"
		row := dbase.QueryRow(sqlStatement, studentId)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		}

		registration := registration{}
		err = row.Scan(&registration.CoursesId, &registration.StudentId, &registration.RegDate)
		if err != nil {
			return c.JSON(http.StatusNotFound, "")
		}

		return c.JSON(http.StatusCreated, registration)
	})
	

	e.POST("/course", func(c echo.Context) error {
		u := new(course)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlStatement := "INSERT INTO public.\"course\" (course_id, course_name, course_fee)VALUES ($1, $2, $3)"
		res, err := dbase.Query(sqlStatement, u.CourseId, u.CourseName, u.CourseFee)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, "ok")
	})

	e.POST("/reg", func(c echo.Context) error {
		u := new(registration)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlStatement := "INSERT INTO public.\"reg\" (courses_id, stud_id, reg_date)VALUES ($1, $2, $3)"
		res, err := dbase.Query(sqlStatement, u.CoursesId, u.StudentId, u.RegDate)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, "ok")
	})

	e.PUT("/course/:courseid", func(c echo.Context) error {
		courseid := c.Param("courseid")
		u := new(course)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlStatement := "UPDATE course SET course_name=$1, course_fee=$2 WHERE course_id=$3"
		res, err := dbase.Query(sqlStatement, u.CourseName, u.CourseFee, courseid)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, u.CourseId)
	})

	e.PUT("/course/:courseid/student/:studentid", func(c echo.Context) error {
		courseid := c.Param("courseid")
		studentid := c.Param("stdentid")

		u := new(registration)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlStatement := "UPDATE reg SET reg_date=$1, courses_id=$2 WHERE courses_id=$2 AND stud_id=$3"
		res, err := dbase.Query(sqlStatement, u.RegDate, courseid, studentid)
		if err != nil {
			fmt.Println(err)
			//return c.JSON(http.StatusCreated, u);
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, registration)
	})

	e.Start(":8080")

}
