You're using a system-managed Python. Use a virtual environment instead:

python3 -m venv venv          # create virtual env
source venv/bin/activate      # activate it each time !!!
pip install pytest            # now install freely
pytest                        # run tests

Your terminal prompt will show (venv) when active. To deactivate later:

deactivate

You'll need to source venv/bin/activate each time you open a new terminal session in that project.
