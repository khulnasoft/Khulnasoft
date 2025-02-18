import pytest
import pytest_asyncio
import os

from logging import warning

from khulnasoft import Sandbox, AsyncSandbox


@pytest.fixture()
def template():
    return "base"


@pytest.fixture()
def sandbox(template, debug):
    sandbox = Sandbox(template)

    try:
        yield sandbox
    finally:
        try:
            sandbox.kill()
        except:
            if not debug:
                warning(
                    "Failed to kill sandbox — this is expected if the test runs with local envd."
                )


@pytest_asyncio.fixture
async def async_sandbox(template, debug):
    sandbox = await AsyncSandbox.create(template)

    try:
        yield sandbox
    finally:
        try:
            await sandbox.kill()
        except:
            if not debug:
                warning(
                    "Failed to kill sandbox — this is expected if the test runs with local envd."
                )


@pytest.fixture
def debug():
    return os.getenv("KHULNASOFT_DEBUG") is not None


@pytest.fixture(autouse=True)
def skip_by_debug(request, debug):
    if request.node.get_closest_marker("skip_debug"):
        if debug:
            pytest.skip("skipped because KHULNASOFT_DEBUG is set")
