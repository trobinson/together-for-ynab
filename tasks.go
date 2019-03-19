package main

type Task struct {
	SourceBudgets []Budget `mapstructure:"source-budgets"`
	JointBudgets  []Budget `mapstructure:"joint-budgets"`
}
