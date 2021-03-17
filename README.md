# App Engine examples
This repo holds examples that I presented in the March 2021 Go Developer's Network (GDN) via GoBridge.

The video is [here](https://youtu.be/YPy6k9O_q1k) and the slides are [here](https://github.com/matt4biz/go-class-slides/blob/trunk/gobridge/gdn-2103-gae-slides.pdf).

## Default
The default app is just a placeholder for the default App Engine app in a project.

## Todo
The `todo` app just takes data from a JSON test site and reformats it into a web page via a template.

The code shows connecting to GCP Stackdriver profiling and has an error that leaks goroutines. It also shows getting the incoming trace tag from a request header and putting that into logs from the app.

## Sort
The `sort` app animates various sorting algorithms (see the code).

The code shows creating and registering custom metrics with OpenCensus which are then exported to Stackdriver.

You can also see the profiling and log trace examples from the `todo` app also.

## Working with the examples
You'll need a GCP project of your own where you can run App Engine apps and use monitoring and profiling.

Each application has its own YAML file for deploying to GCP, e.g., `gcloud app deploy --appyaml=default-app.yaml`.
