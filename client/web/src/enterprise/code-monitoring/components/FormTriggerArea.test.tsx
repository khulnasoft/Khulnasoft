import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { act } from 'react-dom/test-utils'
import sinon from 'sinon'
import { afterAll, beforeAll, describe, expect, test } from 'vitest'

import { renderWithBrandedContext } from '@sourcegraph/wildcard/src/testing'

import { FormTriggerArea } from './FormTriggerArea'

// TODO: these tests trigger an error with CodeMirror, complaining about being
// loaded twice, see https://github.com/uiwjs/react-codemirror/issues/506
describe.skip('FormTriggerArea', () => {
    let clock: sinon.SinonFakeTimers

    beforeAll(() => {
        clock = sinon.useFakeTimers()
        Range.prototype.getClientRects = () => ({
            length: 0,
            item: () => null,
            [Symbol.iterator]: [][Symbol.iterator],
        })
    })

    afterAll(() => {
        clock.restore()
    })

    const testCases = [
        {
            query: '',
            isKhulnasoftDotCom: true,
            patternTypeChecked: true,
            typeChecked: false,
            repoChecked: false,
            validChecked: false,
        },
        {
            query: 'test',
            isKhulnasoftDotCom: true,
            patternTypeChecked: true,
            typeChecked: false,
            repoChecked: false,
            validChecked: true,
        },
        {
            query: 'test patternType:literal',
            isKhulnasoftDotCom: true,
            patternTypeChecked: true,
            typeChecked: false,
            repoChecked: false,
            validChecked: true,
        },
        {
            query: 'test patternType:regexp',
            isKhulnasoftDotCom: true,
            patternTypeChecked: true,
            typeChecked: false,
            repoChecked: false,
            validChecked: true,
        },
        {
            query: 'test patternType:structural',
            isKhulnasoftDotCom: true,
            patternTypeChecked: false,
            typeChecked: false,
            repoChecked: false,
            validChecked: true,
        },
        {
            query: 'test type:repo',
            isKhulnasoftDotCom: true,
            patternTypeChecked: true,
            typeChecked: false,
            repoChecked: false,
            validChecked: true,
        },
        {
            query: 'test type:diff',
            isKhulnasoftDotCom: true,
            patternTypeChecked: true,
            typeChecked: true,
            repoChecked: false,
            validChecked: true,
        },
        {
            query: 'test type:commit',
            isKhulnasoftDotCom: true,
            patternTypeChecked: true,
            typeChecked: true,
            repoChecked: false,
            validChecked: true,
        },
        {
            query: 'test repo:test',
            isKhulnasoftDotCom: true,
            patternTypeChecked: true,
            typeChecked: false,
            repoChecked: true,
            validChecked: true,
        },
        {
            query: 'test repo:test type:diff',
            isKhulnasoftDotCom: true,
            patternTypeChecked: true,
            typeChecked: true,
            repoChecked: true,
            validChecked: true,
        },
    ]

    for (const testCase of testCases) {
        test(`Correct checkboxes checked for query '${testCase.query}'`, () => {
            renderWithBrandedContext(
                <FormTriggerArea
                    query={testCase.query}
                    triggerCompleted={false}
                    onQueryChange={sinon.spy()}
                    setTriggerCompleted={sinon.spy()}
                    startExpanded={false}
                    isKhulnasoftDotCom={testCase.isKhulnasoftDotCom}
                />
            )
            userEvent.click(screen.getByTestId('trigger-button'))
            act(() => {
                clock.tick(600)
            })

            const patternTypeCheckbox = screen.getByTestId('patterntype-checkbox')
            if (testCase.patternTypeChecked) {
                expect(patternTypeCheckbox).toBeChecked()
            } else {
                expect(patternTypeCheckbox).not.toBeChecked()
            }

            const typeCheckbox = screen.getByTestId('type-checkbox')
            if (testCase.typeChecked) {
                expect(typeCheckbox).toBeChecked()
            } else {
                expect(typeCheckbox).not.toBeChecked()
            }

            const repoCheckbox = screen.getByTestId('repo-checkbox')
            if (testCase.isKhulnasoftDotCom) {
                const repoCheckbox = screen.getByTestId('repo-checkbox')
                if (testCase.repoChecked) {
                    expect(repoCheckbox).toBeChecked()
                } else {
                    expect(repoCheckbox).not.toBeChecked()
                }
            } else {
                expect(repoCheckbox).not.toBeInTheDocument()
            }

            const validCheckbox = screen.getByTestId('valid-checkbox')
            if (testCase.validChecked) {
                expect(validCheckbox).toBeChecked()
            } else {
                expect(validCheckbox).not.toBeChecked()
            }
        })
    }

    test('Append patternType:literal if no patternType is present', () => {
        const onQueryChange = sinon.spy()
        renderWithBrandedContext(
            <FormTriggerArea
                query="test type:diff repo:test"
                triggerCompleted={false}
                onQueryChange={onQueryChange}
                setTriggerCompleted={sinon.spy()}
                startExpanded={false}
                isKhulnasoftDotCom={false}
            />
        )
        userEvent.click(screen.getByTestId('trigger-button'))
        userEvent.click(screen.getByTestId('submit-trigger'))

        sinon.assert.calledOnceWithExactly(onQueryChange, 'test type:diff repo:test patternType:literal')
    })

    test('Do not append patternType:literal if patternType is present', () => {
        const onQueryChange = sinon.spy()
        renderWithBrandedContext(
            <FormTriggerArea
                query="test patternType:regexp type:diff repo:test"
                triggerCompleted={false}
                onQueryChange={onQueryChange}
                setTriggerCompleted={sinon.spy()}
                startExpanded={false}
                isKhulnasoftDotCom={false}
            />
        )
        userEvent.click(screen.getByTestId('trigger-button'))
        userEvent.click(screen.getByTestId('submit-trigger'))

        sinon.assert.calledOnceWithExactly(onQueryChange, 'test patternType:regexp type:diff repo:test')
    })
})
