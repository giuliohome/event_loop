import os
import time
import asyncio
import multiprocessing
from dataclasses import dataclass
from datetime import timedelta

from temporalio import activity, workflow
from temporalio.client import Client
from temporalio.worker import SharedStateManager,Worker
from concurrent.futures import ProcessPoolExecutor


# While we could use multiple parameters in the activity, Temporal strongly
# encourages using a single dataclass instead which can have fields added to it
# in a backwards-compatible way.
@dataclass
class ComposeGreetingInput:
    greeting: str
    name: str


# Basic activity that logs and does string concatenation
@activity.defn
def compose_greeting(input: ComposeGreetingInput) -> str:
    activity.logger.info("Running activity with parameter %s" % input)
    myvar = os.getenv('MYVAR')
    os.environ['MYVAR'] = myvar+'!'
    myvar = os.environ['MYVAR']
    time.sleep(5)
    return f"{input.greeting}, {input.name}! My var is {myvar}"


# Basic workflow that logs and invokes an activity
@workflow.defn
class GreetingWorkflow:
    @workflow.run
    async def run(self, name: str) -> str:
        workflow.logger.info("Running workflow with parameter %s" % name)
        task1 = asyncio.create_task( workflow.execute_activity(
            compose_greeting,
            ComposeGreetingInput("Hello", name),
            start_to_close_timeout=timedelta(seconds=10),
        ))
        task2 = asyncio.create_task( workflow.execute_activity(
            compose_greeting,
            ComposeGreetingInput("Hello", name),
            start_to_close_timeout=timedelta(seconds=10),
        ))
        res1 = await task1
        res2 = await task2
        return res1 + " and then " + res2


async def main():
    # Uncomment the lines below to see logging output
    # import logging
    # logging.basicConfig(level=logging.INFO)

    # Start client
    client = await Client.connect("localhost:7233")

    # Run a worker for the workflow
    async with Worker(
        client,
        task_queue="hello-activity-task-queue",
        workflows=[GreetingWorkflow],
        activities=[compose_greeting],
        # Synchronous activities are not allowed unless we provide some kind of
        # executor. Here we are giving a process pool executor which means the
        # activity will actually run in a separate process. This same executor
        # could be passed to multiple workers if desired.
        activity_executor=ProcessPoolExecutor(5),
        # Since we are using an executor that is not a thread pool executor,
        # Temporal needs some kind of manager to share state such as
        # cancellation info and heartbeat info between the host and the
        # activity. Therefore, we must provide a shared_state_manager here. A
        # helper is provided to create it from a multiprocessing manager.
        shared_state_manager=SharedStateManager.create_from_multiprocessing(
            multiprocessing.Manager()
        ),
    ):

        # While the worker is running, use the client to run the workflow and
        # print out its result. Note, in many production setups, the client
        # would be in a completely separate process from the worker.
        result = await client.execute_workflow(
            GreetingWorkflow.run,
            "World",
            id="hello-activity-workflow-id",
            task_queue="hello-activity-task-queue",
        )
        print(f"Result: {result}")


if __name__ == "__main__":
    asyncio.run(main())
