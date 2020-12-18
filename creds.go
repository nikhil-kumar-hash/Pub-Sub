package main

type Creds struct {
	AccessKey string
	SecretKey string
	Region    string
}

func GetCredentials() Creds {
	return Creds{
		AccessKey: "AKIATYO76VW7LIH25VP3",                     //Load your access key here
		SecretKey: "U8LvZDHDa9USFEDgsT4QqTY0rHzfgPQ51G5jcCs9", //Load your secret key here
		Region:    "us-east-2",                                //Load your region here
	}
}
