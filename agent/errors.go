package main

func HandelError(err error, t string){
	switch t {
		case "administrative":
		panic(err)
		case "conntrack":
		println("A problem occured with conntrack connection.")
		println(err)
	}
}
