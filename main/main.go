package main

import (
	"errors"
	"fmt"
	"promise"
)

func main() {

	/*Chain Example1-Then/Catch/Finally*/
	var promise1 = promise.CreatePromise(func(resolve func(interface{}), reject func(error)) {

			const hello = "hello"

			//success
			if hello == "hello" {
				resolve(hello)
				return
			}

			//fake error
			if hello != "hello" {
				reject(errors.New("hello != hello"))
				return
			}
		}).
		Then(func(data interface{}) interface{} {
			fmt.Println("The result is:", data)
			return data.(string) + " world "
		}).
		Then(func(data interface{}) interface{} {
			fmt.Println("The new result is:", data)
			return nil
		}).
		Catch(func(error error) error {
			fmt.Println("Error during execution:", error.Error())
			return nil
		}).
		Finally(func(data interface{}) interface{} {
			fmt.Println("This has to be done")
			return nil
		})

	promise1.Await()

	/*Chain Example2-Multiple Then/Catch*/
	var promise2 = promise.CreatePromise(func(resolve func(interface{}), reject func(error)) {

			const hello2 = "hello23"

			//success
			if hello2 == "hello2" {
				resolve(hello2)
				return
			}

			//fake error
			if hello2 != "hello2" {
				reject(errors.New("hello2 != hello2"))
				return
			}
		}).
		Then(func(data interface{}) interface{} {
			fmt.Println("The result is:", data)
			return data.(string) + " world "
		}).
		Then(func(data interface{}) interface{} {
			fmt.Println("The new result is:", data)
			return nil
		}).
		Catch(func(error error) error {
			fmt.Println("Error during execution:", error.Error()) //This will catch the rejection till this point
			return nil
		}).
		Then(func(data interface{}) interface{} {
			fmt.Println("The result is:", data)
			return data.(string) + " world " //data has become nil so this has to be caught in the next catch
		}).
		Then(func(data interface{}) interface{} {
			fmt.Println("The new result is:", data)
			return nil
		}).
		Catch(func(error error) error {
			fmt.Println("Error during execution:", error.Error())
			return nil
		}).
		Finally(func(data interface{}) interface{} {
			fmt.Println("This has to be done")
			return nil
		})

	promise2.Await()


	/*Chain Example3-Then/Finally*/
	var promise3 = promise.Resolve(nil).
		Then(func(data interface{}) interface{} {
			fmt.Println("1")
			return nil
		}).
		Then(func(data interface{}) interface{} {
			fmt.Println("2")
			return nil
		}).
		Then(func(data interface{}) interface{} {
			fmt.Println("3")
			return nil
		}).
		Finally(func(data interface{}) interface{} {
			fmt.Println("Finally")
			return nil
		})

	promise3.Await()


	/*Basic Resolve*/
	var promise4 = promise.Resolve("Resolve Me")
	result, _ := promise4.Await()
	fmt.Println(result)

	/*Basic Reject*/
	var promise5 = promise.Reject(errors.New("Am an Error"))
	_, err := promise5.Await()
	fmt.Println(err)

}
