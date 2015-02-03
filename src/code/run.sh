#!/bin/bash
cd ./resources/tagAlgo
javac -cp .:./lib/mysql-connector-java-5.1.34-bin.jar Recommendation.java
java -cp .:./lib/mysql-connector-java-5.1.34-bin.jar Recommendation $@
