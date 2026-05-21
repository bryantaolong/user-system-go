package com.bryan.system

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.ConfigurationPropertiesScan
import org.springframework.boot.runApplication

@SpringBootApplication
@ConfigurationPropertiesScan
class UserSystemApplication

fun main(args: Array<String>) {
    runApplication<UserSystemApplication>(*args)
}
