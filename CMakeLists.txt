cmake_minimum_required(VERSION 3.10)
project(my_go_project)

set(ROOT ${CMAKE_CURRENT_SOURCE_DIR})
set(CMD "${ROOT}/cmd")
set(BIN "${ROOT}/bin")

find_program(GO_EXECUTABLE go)
if(NOT GO_EXECUTABLE)
    message(FATAL_ERROR "Go Compiler Not Found")
endif()

file(GLOB TARGETS "${CMD}/*")

set(GO_SOURCE_DIR ${CMAKE_CURRENT_SOURCE_DIR}/pkgs)

file(GLOB_RECURSE GO_FILES ${GO_SOURCE_DIR}/*.go)

foreach(target ${TARGETS})
    get_filename_component(target_name ${target} NAME)
    set(target_output "${BIN}/${target_name}")

    add_custom_command(
        OUTPUT ${target_output}
        COMMAND ${GO_EXECUTABLE} build -o ${target_output} ${CMD}/${target_name}/${target_name}.go
        DEPENDS ${CMD}/${target_name}/${target_name}.go ${GO_FILES}
        WORKING_DIRECTORY ${ROOT}
        COMMENT "Building ${target_name}"
    )

    add_custom_target(${target_name} ALL DEPENDS ${target_output})
    list(APPEND TARGET_OUTPUTS ${target_output})
endforeach()
