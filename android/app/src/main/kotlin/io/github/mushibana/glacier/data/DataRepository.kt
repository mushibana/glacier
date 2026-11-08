package io.github.mushibana.glacier.data

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow


interface DataRepository {
    val data: Flow<List<Task>>
}

data class Task(val id: Int, val name: String, val orderNum: Long = 1024L)

class DefaultDataRepository : DataRepository {
    override val data: Flow<List<Task>> = flow {
        emit(
            (1..25).map { i -> Task(i, "$i task", i * 1024L) }
        )
    }
}
