from khulnasoft import AsyncSandbox


async def test_write_file(async_sandbox: AsyncSandbox):
    filename = "test_write.txt"
    content = "This is a test file."

    info = await async_sandbox.files.write(filename, content)
    assert info.path == f"/home/user/{filename}"

    exists = await async_sandbox.files.exists(filename)
    assert exists

    read_content = await async_sandbox.files.read(filename)
    assert read_content == content


async def test_overwrite_file(async_sandbox: AsyncSandbox):
    filename = "test_overwrite.txt"
    initial_content = "Initial content."
    new_content = "New content."

    await async_sandbox.files.write(filename, initial_content)
    await async_sandbox.files.write(filename, new_content)
    read_content = await async_sandbox.files.read(filename)
    assert read_content == new_content


async def test_write_to_non_existing_directory(async_sandbox: AsyncSandbox):
    filename = "non_existing_dir/test_write.txt"
    content = "This should succeed too."

    await async_sandbox.files.write(filename, content)
    exists = await async_sandbox.files.exists(filename)
    assert exists

    read_content = await async_sandbox.files.read(filename)
    assert read_content == content
