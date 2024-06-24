These are instructions for creating a virtual environment for pytorch
installing pytorch and other requirements, and executing in that venv.

Issues: https://github.com/omry/omegaconf/issues/1107
https://github.com/omry/omegaconf/issues/758

https://github.com/abdeladim-s/easymms/blob/main/README.md

> cd $HOME

> mkdir python_env

> cd python_env

> python3 -m venv pytorch_venv

> source pytorch_venv/bin/activate

Did not update pip despite request, using pip 21.2.4

> pip install hydra-core==1.0.7 # needed to fix fairseq bug

> pip3 install torch torchvision torchaudio

> pip3 install easymms

> pip3 install tensorboardX

> pip3 install "numpy<2"

> deactivate

> export PYTORCH_EXE=$HOME/python_env/pytorch_venv/bin/python3

