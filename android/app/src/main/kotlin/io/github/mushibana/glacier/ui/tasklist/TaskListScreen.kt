package io.github.mushibana.glacier.ui.tasklist

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.consumeWindowInsets
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyListItemInfo
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.rounded.DragIndicator
import androidx.compose.material3.Button
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.Surface
import androidx.compose.material3.Icon
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.text.style.TextAlign
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation3.runtime.NavKey
import io.github.mushibana.glacier.R
import io.github.mushibana.glacier.data.DefaultDataRepository
import io.github.mushibana.glacier.data.Task
import io.github.mushibana.glacier.ui.theme.GlacierTheme
import sh.calvin.reorderable.ReorderableItem
import sh.calvin.reorderable.rememberReorderableLazyListState
import sh.calvin.reorderable.ReorderableCollectionItemScope

@Composable
fun TaskListScreen(
    onItemClick: (NavKey) -> Unit,
    modifier: Modifier = Modifier,
    viewModel: TaskListScreenViewModel = viewModel { TaskListScreenViewModel(DefaultDataRepository()) },
) {
    val uiState by viewModel.uiState.collectAsStateWithLifecycle()
    when (val state = uiState) {
        TaskListScreenUiState.Loading -> {
            Box(
                modifier = modifier.fillMaxSize(),
                contentAlignment = Alignment.Center
            ) {
                CircularProgressIndicator()
            }
        }

        is TaskListScreenUiState.Success -> {
            TaskListScreen(list = state.data, modifier = modifier)
        }

        is TaskListScreenUiState.Error -> {
            Text("Error loading data: ${state.throwable.message}")
        }
    }
}

private fun reorderTasks(tasks: List<Task>, from: LazyListItemInfo, to: LazyListItemInfo): List<Task> {
    if (from.index == to.index) return tasks

    val reordered = tasks.toMutableList().apply {
        add(to.index, removeAt(from.index))
    }

    val beforeOrderNum = reordered.getOrNull(to.index - 1)?.orderNum ?: 0L
    val afterOrderNum = reordered.getOrNull(to.index + 1)?.orderNum ?: (beforeOrderNum + 2048)
    val newOrderNum = beforeOrderNum + (afterOrderNum - beforeOrderNum) / 2

    return if (newOrderNum in (beforeOrderNum + 1) until afterOrderNum) {
        reordered.mapIndexed { index, task ->
            if (index == to.index) {
                task.copy(orderNum = newOrderNum)
            } else {
                task
            }
        }
    } else {
        reordered.mapIndexed { index, task ->
            task.copy(orderNum = (1 + index) * 1024L)
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
internal fun TaskListScreen(list: List<Task>, modifier: Modifier = Modifier) {
    var tasks by remember { mutableStateOf(list) }

    val lazyListState = rememberLazyListState()
    val reorderableLazyListState =
        rememberReorderableLazyListState(lazyListState) { from, to -> tasks = reorderTasks(tasks, from, to) }

    Scaffold(
        modifier = modifier,
        bottomBar = {
            NavigationBar(
                modifier = Modifier.height(80.dp),
                containerColor = MaterialTheme.colorScheme.primaryContainer,
                contentColor = MaterialTheme.colorScheme.primary,
            ) {}
        }
    ) { innerPadding ->
        Surface(
            Modifier
                .fillMaxSize()
                .padding(innerPadding)
                .consumeWindowInsets(innerPadding)
                .padding(top = 20.dp, start = 10.dp, end = 30.dp),
            color = MaterialTheme.colorScheme.surface,
        ) {
            LazyColumn(
                state = lazyListState,
                // verticalArrangement = Arrangement.spacedBy(2.dp),
            ) {
                items(items = tasks, key = { it.id }) { task ->
                    ReorderableItem(state = reorderableLazyListState, key = task.id) { isDragging ->
                        TaskItem(task, isDragging, scope = this)
                    }
                }
            }
        }
    }
}

@Composable
fun TaskItem(task: Task, isDragging: Boolean, scope: ReorderableCollectionItemScope) {
    Row(
        Modifier.fillMaxWidth(),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        IconButton(
            modifier = with(scope) {
                //TODO: return longPress
                // Modifier.longPressDraggableHandle()
                Modifier.draggableHandle()
            },
            onClick = {},
        ) {
            Icon(Icons.Rounded.DragIndicator, contentDescription = "Reorder")
        }
        Text(
            // task.name,
            "${task.orderNum} - ${task.name}",
            textAlign = TextAlign.Start,
            overflow = TextOverflow.Ellipsis,
            maxLines = 1,
            style = MaterialTheme.typography.bodyLarge,
        )
    }
}

@Preview(showBackground = true, showSystemUi = true)
@Composable
fun TaskListScreenPreview() {
    GlacierTheme { TaskListScreen(listOf(Task(1, "Android", 1))) }
}

