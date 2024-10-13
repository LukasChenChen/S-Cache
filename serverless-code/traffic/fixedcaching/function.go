package main

var(
    web_image = "cocc/web:v1.0"
	resizer_image = "cocc/resizer:v1.2"
	checkout_image = "cocc/checkout:v1.3"
	recognition_image = "cocc/yolo:v1.1" //yolo
	container_port  = 3333

	hello_world = Function{
		ApiVersion:     "serving.knative.dev/v1",
		Scheduler:      "fc-scheduler",
		Image:          "gcr.io/knative-samples/helloworld-go",
        Port:           8080,
		Size:           500,
		Type:           1,
		ColdStartTime:  100,
		ProcessTime:    50,
		Clock:          0,
		Priority:       1,
	}
	hello_worlddelete = Function{
		ApiVersion:     "serving.knative.dev/v1",
		Scheduler:      "fc-scheduler",
		Image:          "gcr.io/knative-samples/helloworld-go",
        Port:           8080,
		Size:           500,
		Type:           2,
		ProcessTime:    100,
		Clock:          0,
		Priority:       1,
	}

	container_1 = Function{
		ApiVersion:     "serving.knative.dev/v1",
		Scheduler:      "fc-scheduler",
		Image:          web_image,
        Port:           container_port,
		Size:           55.85,
		Type:           1,
		ColdStartTime:  5.32,
		ProcessTime:    0.055,
		Clock:          0,
		Priority:       1,
		LifeTime:       10,
	}

	container_2 = Function{
		ApiVersion:     "serving.knative.dev/v1",
		Scheduler:      "fc-scheduler",
		Image:          resizer_image,
        Port:           container_port,
		Size:           158.94,
		Type:           2,
		ColdStartTime:  5.33,
		ProcessTime:    0.063,
		Clock:          0,
		Priority:       1,
		LifeTime:       10,
	}

	container_3 = Function{
		ApiVersion:     "serving.knative.dev/v1",
		Scheduler:      "fc-scheduler",
		Image:          checkout_image,
        Port:           container_port,
		Size:           331.11,
		Type:           3,
		ColdStartTime:  5.34,   //[1.7,4.5]
		ProcessTime:    0.202,  //[2,6.5]
		Clock:          0,
		Priority:       1,
		LifeTime:       10,
	}

	container_4 = Function{
		ApiVersion:     "serving.knative.dev/v1",
		Scheduler:      "fc-scheduler",
		Image:          recognition_image,
        Port:           container_port,
		Size:           91.77,
		Type:           4,
		ColdStartTime:  4.89,   //[1.7,4.5]
		ProcessTime:    2.08,  //[2,6.5]
		Clock:          0,
		Priority:       1,
		LifeTime:       10,
	}

	// container_5 = Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           250,
	// 	Type:           5,
	// 	ColdStartTime:  4,   //[1.7,4.5]
	// 	ProcessTime:    3,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_6 = Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           300,
	// 	Type:           6,
	// 	ColdStartTime:  4.5,   //[1.7,4.5]
	// 	ProcessTime:    3.5,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_7 = Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           250,     //[50,300]
	// 	Type:           7,
	// 	ColdStartTime:  4,   //[1.7,4.5]
	// 	ProcessTime:    4,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_8 = Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           200,     //[50,300]
	// 	Type:           8,
	// 	ColdStartTime:  3.5,   //[1.7,4.5]
	// 	ProcessTime:    4.5,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_9 = Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           250,     //[50,300]
	// 	Type:           9,
	// 	ColdStartTime:  3,   //[1.7,4.5]
	// 	ProcessTime:    5,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_10 = Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           200,     //[50,300]
	// 	Type:           10,
	// 	ColdStartTime:  2.5,   //[1.7,4.5]
	// 	ProcessTime:    5.5,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_11 = Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           250,     //[50,300]
	// 	Type:           11,
	// 	ColdStartTime:  2,   //[1.7,4.5]
	// 	ProcessTime:    6,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_12= Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           200,     //[50,300]
	// 	Type:           12,
	// 	ColdStartTime:  1.7,   //[1.7,4.5]
	// 	ProcessTime:    6.5,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_13= Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           250,     //[50,300]
	// 	Type:           13,
	// 	ColdStartTime:  2,   //[1.7,4.5]
	// 	ProcessTime:    6,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_14= Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           300,     //[50,300]
	// 	Type:           14,
	// 	ColdStartTime:  2.5,   //[1.7,4.5]
	// 	ProcessTime:    5.5,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_15= Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           300,     //[50,300]
	// 	Type:           15,
	// 	ColdStartTime:  3.5,   //[1.7,4.5]
	// 	ProcessTime:    4.5,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_16= Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           250,     //[50,300]
	// 	Type:           16,
	// 	ColdStartTime:  4,   //[1.7,4.5]
	// 	ProcessTime:    4,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_17= Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           200,     //[50,300]
	// 	Type:           17,
	// 	ColdStartTime:  4.5,   //[1.7,4.5]
	// 	ProcessTime:    3.5,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_18= Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           250,     //[50,300]
	// 	Type:           18,
	// 	ColdStartTime:  4,   //[1.7,4.5]
	// 	ProcessTime:    3,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_19= Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           300,     //[50,300]
	// 	Type:           19,
	// 	ColdStartTime:  3.5,   //[1.7,4.5]
	// 	ProcessTime:    2.5,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }

	// container_20= Function{
	// 	ApiVersion:     "serving.knative.dev/v1",
	// 	Scheduler:      "my-scheduler",
	// 	Image:          container_image,
    //     Port:           container_port,
	// 	Size:           250,     //[50,300]
	// 	Type:           20,
	// 	ColdStartTime:  3,   //[1.7,4.5]
	// 	ProcessTime:    2,  //[2,6.5]
	// 	Clock:          0,
	// 	Priority:       1,
	// }
)

