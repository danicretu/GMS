#!/bin/bash -p
clear
javac -cp lib/mysql-connector-java-5.1.34-bin.jar src/Recommendation.java
java -cp .:lib/mysql-connector-java-5.1.34-bin.jar:./src Recommendation $@
