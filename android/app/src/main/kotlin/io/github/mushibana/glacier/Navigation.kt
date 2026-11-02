package io.github.mushibana.glacier

import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.safeDrawingPadding
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation3.runtime.entryProvider
import androidx.navigation3.runtime.rememberNavBackStack
import androidx.navigation3.ui.NavDisplay
import io.github.mushibana.glacier.ui.tasklist.TaskListScreen

@Composable
fun MainNavigation() {
  val backStack = rememberNavBackStack(TaskList)

  NavDisplay(
    backStack = backStack,
    onBack = { backStack.removeLastOrNull() },
    entryProvider =
      entryProvider {
        entry<TaskList> {
          TaskListScreen(onItemClick = { navKey -> backStack.add(navKey) })
        }
      },
  )
}
