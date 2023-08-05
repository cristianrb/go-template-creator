run:
	cd cmd && go run main.go gtc --routing=chi --logging=zap --config=viper --path=/home/cristian/development/go --name=gtc_test

run_short:
	cd cmd && go run main.go gtc -r=chi -l=zap -c=viper

run_default:
	cd cmd && go run main.go gtc

help:
	cd cmd && go run main.go gtc --help