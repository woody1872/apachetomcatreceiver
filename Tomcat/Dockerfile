FROM tomcat:latest

RUN mkdir -p /usr/local/tomcat/webapps/testapp

COPY index.html /usr/local/tomcat/webapps/testapp/index.html

COPY tomcat-users.xml /usr/local/tomcat/conf/tomcat-users.xml

ADD manager/ /usr/local/tomcat/webapps/manager/