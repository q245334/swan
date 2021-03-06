<!--
 Copyright (c) 2017 Intel Corporation

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
-->

# Jupyter experiment viewer

## Introduction

> "The Jupyter Notebook is a web application that allows you to create and share documents that contain live code, equations, visualizations and explanatory text. Uses include: data cleaning and transformation, numerical simulation, statistical modeling, machine learning and much more." [from jupyter.org](http://jupyter.org/)

Swan uses *Jupyter Notebook* to filter, process and visualize results from experiments.

## Installation

In order to install Jupyter, you need `python` and `pip` installed. On Centos 7, installation of those packages can be achieved by:

```sh
sudo yum install python-pip python-devel
```

or following the instructions at the [official `pip` site](https://pip.pypa.io/en/stable/installing/#installing-with-get-pip-py).

After `python` and `pip` are installed Jupyter can be installed by typing:

```sh
make deps_jupyter
```
in the Swan root directory.

## Launching Jupyter

In order to start Jupyter go to the Jupyter directory in Swan repository then run `notebook`:

```sh
cd jupyter/
jupyter notebook
```

Jupyter will start locally. Point a web browser to http://127.0.0.1:8888 to access Jupyter notebooks.


## Explore the Example Jupyter Notebook

From within the Jupyter web interface, open a template notebook by clicking on `example.ipynb`:

![experiment](/images/new_notebook.png)

This is very simple notebook that will generate only sensitivity profile for the experiment.
The first step is to set the following variables:
- `IP` and `PORT` which will point to the Cassandra database that contains the experiment's results and metadata
- `EXPERIMENT_ID` is the identifier of the experiment which will be examined

After filling the variables, navigate to the green box using the keyboard arrows so that it points to the first variable and press `[Shift] [Enter]` to evaluate it. Evaluation actually means executing the code in the box. Evaluate further and observe the output. `Experiment` object's construction will look like:

```python
# An experiment can now be loaded from the database by its ID.
exp1 = Experiment(cassandra_cluster=[IP],
                  experiment_id=EXPERIMENT_ID_SWAN, port=PORT, name="experiment 1")
```
It may take a while since it will retrieve data from Cassandra and store it in the variable `exp1` which represents itself as a table:

![sample list](/images/sample_list.png)

The last two steps are to render the sensitivity profile from the loaded samples and draw sensitivity chart. The former will be generated after evaluating:

```python
p = Profile(exp1, slo=500)
p.sensitivity_table()
```

Where `slo` is the target latency in microseconds. The rendered table, which is the *Sensitivity Profile*, will look similar to the one below:

![sensitivity profile](/images/sensitivity_profile.png)

To learn more about *Sensitivity Profile* read the [Sensitivity Experiment](experiments/memcached-sensitivity-profile/README.md) README.

## Visualizing data using Jupyter

We will use another tool, called [plotly](https://plot.ly/), to interactively plot our data points. You will need to [install it manually](https://plot.ly/python/getting-started/#installation).

Once installed, render a *Sensitivity Chart* from a profile using the following method:

```python
p1.sensitivity_chart(fill=True, to_max=False)
```

A *Sensitivity Chart* will show how aggressors interference with the High Priority job.

![sensitivity_chart](/images/sensitivity_chart.png)

The chart shows *Sensitivity Profile* in the form of a graph. The red horizontal line is the SLO which cannot be violated. Below the SLO line there is a Baseline and for all load points it should stay there. For Baseline there must not be any SLO violation. High Priority (HP) in some cases will start below the SLO may rapidly exceed the SLO at some load point.

The `fill` parameter to the `sensitivity_chart()` controls whether a color fills the area between the Baseline and aggressors. `to_max` shows a comparison between Baseline and the 'worst case', which for each load point takes the highest latency from all HP with aggressor.


It is also possible to compare two experiments (examine `example_2.ipynb`):

```python
exps = [exp1, exp2]
compare_experiments(exps, fill=True, to_max=False)
```

Here `fill` parameter acts the same as in the previous example, and `to_max` compares baseline for two experiments.

![compare_two_experiments](/images/compare_two_experiments.png)

At this chart the "green area" shows improvement in terms of higher load and lower latency, between `Baselines` on two different setups.

## Exploration data using Jupyter Notebooks

To get started, we have provided example notebooks [example.ipynb](example.ipynb) and [example_2.ipynb](example_2.ipynb).
