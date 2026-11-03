package io.github.mushibana.glacier.data

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow


interface DataRepository {
    val data: Flow<List<Task>>
}

data class Task(val name: String, val orderNum: Int)

class DefaultDataRepository : DataRepository {
    override val data: Flow<List<Task>> = flow {
        emit(
            listOf(
                Task("5th test task", 5 * 1024),
                Task("4th test task", 4 * 1024),
                Task("3rd test task", 3 * 1024),
                Task("2nd test task", 2 * 1024),
                Task("1st test task", 1 * 1024),
            )
        )
    }
}
