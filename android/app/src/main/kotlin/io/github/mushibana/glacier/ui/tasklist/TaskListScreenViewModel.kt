package io.github.mushibana.glacier.ui.main

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import io.github.mushibana.glacier.data.DataRepository
import io.github.mushibana.glacier.ui.tasklist.TaskListScreenUiState.Success
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.flow.stateIn

class TaskListScreenViewModel(dataRepository: DataRepository) : ViewModel() {
  val uiState: StateFlow<TaskListScreenUiState> =
    dataRepository.data
      .map<List<Task>, TaskListScreenUiState>(::Success)
      .catch { emit(TaskListScreenUiState.Error(it)) }
      .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5000), TaskListScreenUiState.Loading)
}

sealed interface TaskListScreenUiState {
  object Loading : TaskListScreenUiState

  data class Error(val throwable: Throwable) : TaskListScreenUiState

  data class Success(val data: List<String>) : TaskListScreenUiState
}
