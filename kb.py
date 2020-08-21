import itertools
import time

from utils.aima_utils import (
    Expr, expr
)

from utils.aima_logic import (
    FolKB, pl_resolution, standardize_variables, parse_definite_clause,
    unify, subst, is_var_symbol, constant_symbols, variables
)


def gomas_subst(s, x):
    if isinstance(x, list):
        return [gomas_subst(s, xi) for xi in x]
    elif isinstance(x, tuple):
        return tuple([gomas_subst(s, xi) for xi in x])
    elif not isinstance(x, Expr):
        return expr(x)
    elif is_var_symbol(x.op):
        v = s.get(x.op, x)
        if not isinstance(v, Expr):
            return expr(v)
        else:
            return v
    else:
        return Expr(x.op, *[gomas_subst(s, arg) for arg in x.args])


class kb_engine(FolKB):
    def __init__(self):
        super().__init__()
        self.volatile_memory = {}
        self.var_to_facts = {}

    def memorize(self, keys, conditions):
        assert isinstance(keys, list) or isinstance(keys, tuple)
        assert isinstance(conditions, list) or isinstance(conditions, tuple)
        self.var_to_facts[tuple(keys)] = conditions

    def sense(self, key, value):
        self.volatile_memory[key] = value

    def tell(self, expression):
        if isinstance(expression, list) or isinstance(expression, tuple):
            for e in expressions:
                self.tell(e)
        elif isinstance(expression, str):
            expression = expr(expression)
        assert isinstance(expression, Expr)
        super().tell(expression)

    def ask(self, statement):
        if isinstance(statement, str):
            statement = expr(statement)
        assert isinstance(statement, Expr)
        return super().ask(statement)

    def clear(self):
        self.clauses = []
        self.volatile_memory = {}

    def generate_facts(self, variables):
        evaluations = self.var_to_facts
        permanent_facts = {}
        permanent_facts.update(variables)
        keys = set(permanent_facts)
        generated_facts = []

        for topic in evaluations:
            if keys.intersection(topic):
                for fact, evaluation in evaluations[topic]:
                    try:
                        if eval(evaluation, permanent_facts):
                            generated_facts.append(gomas_subst(permanent_facts, expr(fact)))
                            # generated_facts.append(subst(permanent_facts, expr(fact)))
                    except Exception as ex:
                        if 'Train' in fact:
                            print(str(ex))
                        # raise Exception('Problem in generate_facts: %s' % ('evaluation',))
        return generated_facts

    """
        Check if FolKB entails the query, alpha
    """
    def fol_fc_ask(self, alpha):
        clauses = []
        clauses.extend(self.generate_facts(self.volatile_memory))
        clauses.extend(self.clauses)

        if isinstance(alpha, str):
            alpha = expr(alpha)
        """A simple forward-chaining algorithm. [Figure 9.3]"""
        # TODO: Improve efficiency
        kb_consts = list({c for clause in clauses for c in constant_symbols(clause)})
        def enum_subst(p):
            query_vars = list({v for clause in p for v in variables(clause)})
            for assignment_list in itertools.product(kb_consts, repeat=len(query_vars)):
                theta = {x: y for x, y in zip(query_vars, assignment_list)}
                yield theta

        # check if we can answer without new inferences
        for q in clauses:
            phi = unify(q, alpha, {})
            if phi is not None:
                yield phi

        while True:
            new = []
            for rule in clauses:
                p, q = parse_definite_clause(rule)
                for theta in enum_subst(p):
                    if set(subst(theta, p)).issubset(set(clauses)):
                        q_ = subst(theta, q)
                        if all([unify(x, q_, {}) is None for x in clauses + new]):
                            new.append(q_)
                            phi = unify(q_, alpha, {})

                            if phi is not None:
                                yield phi
            if not new:
                break
            for clause in new:
                clauses.append(clause)
        return None
