module azzadigital.com/tempmeasurement/cloudfunctions

go 1.14

replace azzadigital.com/tempmeasurement/cloudfunctions/dbloader => ./dbloader

require (
	azzadigital.com/tempmeasurement/cloudfunctions/dbloader v0.0.0-00010101000000-000000000000
	cloud.google.com/go v0.63.0
	github.com/GoogleCloudPlatform/functions-framework-go v1.1.0
	google.golang.org/api v0.30.0
)
